package market

import (
	"errors"
	"fmt"
	"scoing-trader/trader/model/market/model"
	"sort"
	"time"
)

type Accountant struct {
	InitialBalance float64
	Fee            float64
	Balance        float64
	Market         model.Market
	Positions      map[string]map[float64]float64
	Assets         map[string]float64
	CoinValues     map[string]float64
}

func NewAccountant(market model.Market, initialBalance float64, fee float64) *Accountant {
	return &Accountant{
		InitialBalance: initialBalance,
		Fee:            fee,
		Balance:        initialBalance,
		Market:         market,
		Positions:      make(map[string]map[float64]float64),
		Assets:         make(map[string]float64),
		CoinValues:     make(map[string]float64),
	}
}

func (w *Accountant) Buy(coin string, quantity float64) error {

	if quantity < 0 {
		return errors.New(fmt.Sprintf("negative buy quantity %2.f", quantity))
	}

	transactionValue := w.CoinValues[coin] * quantity * (1 + w.Fee)

	if transactionValue > w.Balance {
		return errors.New(fmt.Sprintf("buy transaction: %2.f exceeds balance: %.2f", transactionValue, w.Balance))
	}

	buyOrder := model.OrderRequest{
		Symbol:    coin,
		Side:      model.BUY,
		Type:      model.MARKET,
		Timestamp: w.GetTimeStamp(),
		Quantity:  quantity,
	}

	if err := w.Market.NewOrder(buyOrder); err != nil {
		return err
	}

	w.Balance -= transactionValue

	if _, hasKey := w.Positions[coin]; !hasKey {
		w.Positions[coin] = make(map[float64]float64)
	}

	if _, hasKey := w.Positions[coin][w.CoinValues[coin]]; hasKey {
		w.Positions[coin][w.CoinValues[coin]] += quantity
	} else {
		w.Positions[coin][w.CoinValues[coin]] = quantity
	}

	if _, hasKey := w.Assets[coin]; !hasKey {
		w.Assets[coin] = quantity
	} else {
		w.Assets[coin] += quantity
	}

	return nil
}

func (w *Accountant) Sell(coin string, quantity float64) error {
	if quantity < 0 {
		return errors.New(fmt.Sprintf("negative buy quantity %2.f", quantity))
	}

	if quantity > w.Assets[coin] {
		return errors.New(fmt.Sprintf("sell quantity: %2.f exceeds available: %.2f", quantity, w.Assets[coin]))
	}

	sellOrder := model.OrderRequest{
		Symbol:    coin,
		Side:      model.SELL,
		Type:      model.MARKET,
		Timestamp: w.GetTimeStamp(),
		Quantity:  quantity,
	}

	if err := w.Market.NewOrder(sellOrder); err != nil {
		return err
	}

	positionBuyValues := make([]float64, 0, len(w.Positions[coin]))
	for pbv := range w.Positions[coin] {
		positionBuyValues = append(positionBuyValues, pbv)
	}

	sort.Float64s(positionBuyValues)
	remainingQty := quantity

	for _, pbv := range positionBuyValues {
		if remainingQty >= w.Positions[coin][pbv] {
			remainingQty -= w.Positions[coin][pbv]
			delete(w.Positions[coin], pbv)
		} else {
			w.Positions[coin][pbv] -= remainingQty
			remainingQty = 0
		}
		if remainingQty == 0 {
			break
		}
	}

	w.Assets[coin] -= quantity

	currentValue := w.CoinValues[coin]
	w.Balance += currentValue * quantity * (1 - w.Fee)

	return nil
}

func (w *Accountant) UpdateCoinValue(coin string, value float64, timestamp time.Time) error {
	if value < 0 {
		return errors.New(fmt.Sprintf("negative coin value (%s,%.2f)", coin, value))
	}
	w.CoinValues[coin] = value
	w.Market.UpdateCoinValue(coin, value)
	return nil
}

func (w *Accountant) GetBalance() float64 {
	return w.Balance
}

func (w *Accountant) GetPositions(coin string) map[float64]float64 {
	return w.Positions[coin]
}

func (w *Accountant) GetFee() float64 {
	return w.Fee
}

func (w *Accountant) TotalAssetValue() float64 {
	var totalValue float64 = 0
	for coin, assetValue := range w.Assets {
		totalValue += assetValue * w.CoinValues[coin]
	}
	return totalValue
}

func (w *Accountant) NetWorth() float64 {
	return w.Balance + w.TotalAssetValue()
}

func (w *Accountant) CoinPortfolioValue(coin string) float64 {
	if value, contains := w.Assets[coin]; contains {
		return value
	} else {
		return 0.0
	}
}

func (w *Accountant) ToString() string {
	walletStr := fmt.Sprintf(">> NW:%.2f Balance:%.2f |", w.NetWorth(), w.GetBalance())

	var coinList []string

	for coin, _ := range w.CoinValues {
		coinList = append(coinList, coin)
	}

	sort.Strings(coinList)

	for _, coin := range coinList {
		walletStr += fmt.Sprintf(" %s #%d Total:%.2f$(%.2f%%) |", coin, len(w.Positions[coin]), w.CoinPortfolioValue(coin),
			w.CoinPortfolioValue(coin)/w.NetWorth()*100)
	}

	return walletStr
}

func (w *Accountant) GetTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
