package market

import (
	"errors"
	"fmt"
	"math/rand"
	"scoing-trader/trader/model/market/model"
)

type SimulatedMarket struct {
	accountInfo  model.AccountInformation
	orderList    []*model.OrderResponseFull
	tradeList    []*model.Trade
	coinValues   map[string]float64
	unfilledRate float64
	fee          float64
}

func NewSimulatedMarket(unfilledRate float64, fee float64) *SimulatedMarket {
	return &SimulatedMarket{
		accountInfo:  model.AccountInformation{},
		orderList:    make([]*model.OrderResponseFull, 0),
		tradeList:    make([]*model.Trade, 0),
		coinValues:   make(map[string]float64),
		unfilledRate: unfilledRate,
		fee:          fee,
	}
}

func (s *SimulatedMarket) NewOrder(order model.OrderRequest) error {
	if len(s.OpenOrders(order.Symbol)) > 0 {
		return errors.New("Order already open for symbol: " + order.Symbol)
	}

	status := model.NEW
	var execQty float64
	var cumQty float64

	assetBalanceIdx, asset, quoteBalanceIdx, quote := s.getAssetQuoteIdx(order.Symbol)

	if order.Price == 0.0 {
		order.Price = s.coinValues[order.Symbol]
	}

	if quoteBalanceIdx == -1 {
		return errors.New(fmt.Sprintf("No balance of asset: %s ", quote))
	}

	if assetBalanceIdx == -1 {
		s.accountInfo.Balances = append(s.accountInfo.Balances, model.Balance{
			Asset:  asset,
			Free:   0,
			Locked: 0,
		})

		assetBalanceIdx = len(s.accountInfo.Balances) - 1
	}

	if order.Side == model.SELL && s.accountInfo.Balances[assetBalanceIdx].Free < order.Quantity {
		return errors.New(fmt.Sprintf("asset balance for %s insufficent (%f)", asset, order.Quantity))
	}

	if order.Side == model.BUY && s.accountInfo.Balances[quoteBalanceIdx].Free < order.Quantity*order.Price*(1+s.fee) {
		return errors.New(fmt.Sprintf("balance for %s (%f) doesn't cover transaction (%f)", quote,
			s.accountInfo.Balances[quoteBalanceIdx].Free, order.Quantity*order.Price*(1+s.fee)))
	}

	if rand.Float64() >= s.unfilledRate {
		status = model.FILLED
		execQty = order.Quantity
		cumQty = order.QuoteOrderQty

		if order.Side == model.BUY {
			s.accountInfo.Balances[quoteBalanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[quoteBalanceIdx].Free - (order.Quantity * order.Price * (1 + s.fee)))
			s.accountInfo.Balances[assetBalanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[assetBalanceIdx].Free + order.Quantity)
		} else if order.Side == model.SELL {
			s.accountInfo.Balances[assetBalanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[assetBalanceIdx].Free - order.Quantity)
			s.accountInfo.Balances[quoteBalanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[quoteBalanceIdx].Free + (order.Quantity * order.Price * (1 - s.fee)))
		}

		//TODO: discount using BNB as commission asset
		trade := model.Trade{
			Symbol:          order.Symbol,
			Id:              0,
			OrderId:         0,
			OrderListId:     0,
			Price:           order.Price,
			Qty:             order.Quantity,
			Commission:      (order.Price * order.Quantity) * s.fee,
			CommissionAsset: "USDT",
			Time:            0,
			IsBuyer:         order.Side == model.BUY,
			IsMaker:         false,
			IsBestMatch:     true,
		}

		s.tradeList = append(s.tradeList, &trade)
	} else {
		if order.Side == model.BUY {
			s.accountInfo.Balances[quoteBalanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[quoteBalanceIdx].Free - (order.Quantity * order.Price * (1 + s.fee)))
			s.accountInfo.Balances[quoteBalanceIdx].Locked = model.TruncateFloat(s.accountInfo.Balances[quoteBalanceIdx].Locked + (order.Quantity * order.Price * (1 + s.fee)))
		} else if order.Side == model.SELL {
			s.accountInfo.Balances[assetBalanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[assetBalanceIdx].Free - order.Quantity)
			s.accountInfo.Balances[assetBalanceIdx].Locked = model.TruncateFloat(s.accountInfo.Balances[assetBalanceIdx].Locked + order.Quantity)
		}
	}

	orderResp := model.OrderResponseFull{
		Symbol:              order.Symbol,
		OrderId:             -1,
		OrderListId:         -1,
		ClientOrderId:       order.ClientOrderId,
		TransactionTime:     -1,
		Price:               order.Price,
		OrigQty:             order.Quantity,
		ExecutedQty:         execQty,
		CummulativeQuoteQty: cumQty,
		Status:              status,
		TimeInForce:         order.TimeInForce,
		Type:                order.Type,
		Side:                order.Side,
		Fills:               nil,
	}

	s.orderList = append(s.orderList, &orderResp)

	return nil
}

func (s *SimulatedMarket) OpenOrders(symbol string) []*model.OrderResponseFull {
	var openOrders []*model.OrderResponseFull

	for _, order := range s.orderList {
		if order.Status == model.NEW && order.Symbol == symbol {
			openOrders = append(openOrders, order)
		}
	}

	return openOrders
}

func (s *SimulatedMarket) OrderHistory() []*model.OrderResponseFull {
	return s.orderList
}

func (s *SimulatedMarket) CancelOrder(orderId string) error {
	for _, order := range s.orderList {
		if order.ClientOrderId == orderId {
			order.Status = model.CANCELLED

			_, _, quoteBalanceIdx, quote := s.getAssetQuoteIdx(order.Symbol)
			if quoteBalanceIdx == -1 {
				return errors.New(fmt.Sprintf("Could not find balances for quote:%s", quote))
			}

			s.accountInfo.Balances[quoteBalanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[quoteBalanceIdx].Free + s.accountInfo.Balances[quoteBalanceIdx].Locked)
			s.accountInfo.Balances[quoteBalanceIdx].Locked = 0

			return nil
		}
	}
	return errors.New("unknown order id")
}

func (s *SimulatedMarket) AccountInformation() model.AccountInformation {
	return s.accountInfo
}

func (s *SimulatedMarket) Balance(asset string) (model.Balance, error) {
	for _, balance := range s.accountInfo.Balances {
		if balance.Asset == asset {
			return balance, nil
		}
	}
	return model.Balance{}, errors.New("balance for asset " + asset + " does not exist")
}

func (s *SimulatedMarket) Trades() []*model.Trade {
	return s.tradeList
}

func (s *SimulatedMarket) UpdateInformation() {}

func (s *SimulatedMarket) CoinValue(asset string) (float64, error) {
	coinVal, exists := s.coinValues[asset]
	if exists {
		return coinVal, nil
	} else {
		return 0.0, errors.New("asset " + asset + " does not exist")
	}
}

func (s *SimulatedMarket) Deposit(asset string, qty float64) {
	balanceIdx := -1

	for idx, balance := range s.accountInfo.Balances {
		if balance.Asset == asset {
			balanceIdx = idx
		}
	}

	if balanceIdx == -1 {
		s.accountInfo.Balances = append(s.accountInfo.Balances, model.Balance{
			Asset:  asset,
			Free:   qty,
			Locked: 0,
		})
	} else {
		s.accountInfo.Balances[balanceIdx].Free = model.TruncateFloat(s.accountInfo.Balances[balanceIdx].Free + qty)
	}
}

func (s *SimulatedMarket) UpdateCoinValue(asset string, value float64) {
	s.coinValues[asset] = value
}

func (s *SimulatedMarket) getAssetQuoteIdx(symbol string) (int, string, int, string) {
	//TODO: generalize
	asset := symbol[0:3]
	quote := symbol[3:]

	assetBalanceIdx := -1
	quoteBalanceIdx := -1

	for idx := range s.accountInfo.Balances {
		if s.accountInfo.Balances[idx].Asset == asset {
			assetBalanceIdx = idx
		} else if s.accountInfo.Balances[idx].Asset == quote {
			quoteBalanceIdx = idx
		}
	}

	return assetBalanceIdx, asset, quoteBalanceIdx, quote
}
