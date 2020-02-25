package market

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"scoing-trader/trader/model/market/model"
)

type SimulatedMarket struct {
	accountInfo  model.AccountInformation
	orderList    []*model.OrderResponseFull
	tradeList    []*model.Trade
	coinValues   map[string]decimal.Decimal
	unfilledRate float64
	fee          decimal.Decimal
}

func NewSimulatedMarket(unfilledRate float64, fee decimal.Decimal) *SimulatedMarket {
	return &SimulatedMarket{
		accountInfo:  model.AccountInformation{},
		orderList:    make([]*model.OrderResponseFull, 0),
		tradeList:    make([]*model.Trade, 0),
		coinValues:   make(map[string]decimal.Decimal),
		unfilledRate: unfilledRate,
		fee:          fee,
	}
}

func (s *SimulatedMarket) NewOrder(order model.OrderRequest) error {
	if len(s.OpenOrders(order.Symbol)) > 0 {
		return errors.New("Order already open for symbol: " + order.Symbol)
	}

	status := model.NEW
	var execQty decimal.Decimal
	var cumQty decimal.Decimal

	assetBalanceIdx, asset, quoteBalanceIdx, quote := s.getAssetQuoteIdx(order.Symbol)

	if order.Price.IsZero() {
		order.Price = s.coinValues[order.Symbol]
	}

	if quoteBalanceIdx == -1 {
		return errors.New(fmt.Sprintf("No balance of asset: %s ", quote))
	}

	if assetBalanceIdx == -1 {
		s.accountInfo.Balances = append(s.accountInfo.Balances, model.Balance{
			Asset:  asset,
			Free:   decimal.Zero,
			Locked: decimal.Zero,
		})

		assetBalanceIdx = len(s.accountInfo.Balances) - 1
	}

	if order.Side == model.SELL && s.accountInfo.Balances[assetBalanceIdx].Free.LessThan(order.Quantity) {
		return errors.New(fmt.Sprintf("asset balance for %s insufficent (%s)", asset, order.Quantity))
	}

	if order.Side == model.BUY && s.accountInfo.Balances[quoteBalanceIdx].Free.LessThan(order.Quantity.Mul(order.Price).Mul(decimal.NewFromInt(1).Add(s.fee))) {
		return errors.New(fmt.Sprintf("balance for %s (%s) doesn't cover transaction (%s)", quote,
			s.accountInfo.Balances[quoteBalanceIdx].Free, order.Quantity.Mul(order.Price).Mul(decimal.NewFromInt(1).Add(s.fee))))
	}

	if rand.Float64() >= s.unfilledRate {
		status = model.FILLED
		execQty = order.Quantity
		cumQty = order.QuoteOrderQty

		if order.Side == model.BUY {
			s.accountInfo.Balances[quoteBalanceIdx].Free = s.accountInfo.Balances[quoteBalanceIdx].Free.Sub(
				order.Quantity.Mul(order.Price).Mul(decimal.NewFromInt(1).Add(s.fee)))
			s.accountInfo.Balances[assetBalanceIdx].Free = s.accountInfo.Balances[assetBalanceIdx].Free.Add(order.Quantity)
		} else if order.Side == model.SELL {
			s.accountInfo.Balances[assetBalanceIdx].Free = s.accountInfo.Balances[assetBalanceIdx].Free.Sub(order.Quantity)
			s.accountInfo.Balances[quoteBalanceIdx].Free = s.accountInfo.Balances[quoteBalanceIdx].Free.Add(
				order.Quantity.Mul(order.Price).Mul(decimal.NewFromInt(1).Sub(s.fee)))
		}

		//TODO: discount using BNB as commission asset
		trade := model.Trade{
			Symbol:          order.Symbol,
			Id:              0,
			OrderId:         0,
			OrderListId:     0,
			Price:           order.Price,
			Qty:             order.Quantity,
			Commission:      order.Price.Mul(order.Quantity).Mul(s.fee),
			CommissionAsset: "USDT",
			Time:            0,
			IsBuyer:         order.Side == model.BUY,
			IsMaker:         false,
			IsBestMatch:     true,
		}

		s.tradeList = append(s.tradeList, &trade)
	} else {
		if order.Side == model.BUY {
			s.accountInfo.Balances[quoteBalanceIdx].Free = s.accountInfo.Balances[quoteBalanceIdx].Free.Sub(
				order.Quantity.Mul(order.Price).Mul(decimal.NewFromInt(1).Add(s.fee)))
			s.accountInfo.Balances[quoteBalanceIdx].Locked = s.accountInfo.Balances[quoteBalanceIdx].Locked.Add(
				order.Quantity.Mul(order.Price).Mul(decimal.NewFromInt(1).Add(s.fee)))
		} else if order.Side == model.SELL {
			s.accountInfo.Balances[assetBalanceIdx].Free = s.accountInfo.Balances[assetBalanceIdx].Free.Sub(order.Quantity)
			s.accountInfo.Balances[assetBalanceIdx].Locked = s.accountInfo.Balances[assetBalanceIdx].Locked.Add(order.Quantity)
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

			s.accountInfo.Balances[quoteBalanceIdx].Free = s.accountInfo.Balances[quoteBalanceIdx].Free.Add(s.accountInfo.Balances[quoteBalanceIdx].Locked)
			s.accountInfo.Balances[quoteBalanceIdx].Locked = decimal.Zero

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

func (s *SimulatedMarket) CoinValue(asset string) (decimal.Decimal, error) {
	coinVal, exists := s.coinValues[asset]
	if exists {
		return coinVal, nil
	} else {
		return decimal.Zero, errors.New("asset " + asset + " does not exist")
	}
}

func (s *SimulatedMarket) Deposit(asset string, qty decimal.Decimal) {
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
			Locked: decimal.Zero,
		})
	} else {
		s.accountInfo.Balances[balanceIdx].Free = s.accountInfo.Balances[balanceIdx].Free.Add(qty)
	}
}

func (s *SimulatedMarket) UpdateCoinValue(asset string, value decimal.Decimal) {
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
