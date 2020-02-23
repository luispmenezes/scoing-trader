package market

import (
	"errors"
	"fmt"
	"scoing-trader/trader/model/market/model"
	"sort"
	"time"
)

type Accountant struct {
	InitialBalance int64
	Fee            float64
	Balance        int64
	Market         model.Market
	Positions      map[string]map[int64]float64
	Assets         map[string]float64
	AssetValues    map[string]int64
}

func NewAccountant(market model.Market, initialBalance int64, fee float64) *Accountant {
	return &Accountant{
		InitialBalance: initialBalance,
		Fee:            fee,
		Balance:        initialBalance,
		Market:         market,
		Positions:      make(map[string]map[int64]float64),
		Assets:         make(map[string]float64),
		AssetValues:    make(map[string]int64),
	}
}

func (a *Accountant) Buy(coin string, quantity float64) error {

	if quantity < 0 {
		return errors.New(fmt.Sprintf("negative buy quantity %.8f", quantity))
	}

	transactionValue := model.IntFloatMul(a.AssetValues[coin], quantity*(1+a.Fee))

	if transactionValue > a.Balance {
		return errors.New(fmt.Sprintf("buy transaction: %s exceeds balance: %s",
			model.IntToString(transactionValue), model.IntToString(a.Balance)))
	}

	buyOrder := model.OrderRequest{
		Symbol:    coin,
		Side:      model.BUY,
		Type:      model.MARKET,
		Timestamp: a.GetTimeStamp(),
		Quantity:  quantity,
	}

	if err := a.Market.NewOrder(buyOrder); err != nil {
		return err
	}

	a.Balance -= transactionValue

	if _, hasKey := a.Positions[coin]; !hasKey {
		a.Positions[coin] = make(map[int64]float64)
	}

	if _, hasKey := a.Positions[coin][a.AssetValues[coin]]; hasKey {
		a.Positions[coin][a.AssetValues[coin]] += quantity
	} else {
		a.Positions[coin][a.AssetValues[coin]] = quantity
	}

	if _, hasKey := a.Assets[coin]; !hasKey {
		a.Assets[coin] = quantity
	} else {
		a.Assets[coin] += quantity
	}

	return nil
}

func (a *Accountant) Sell(coin string, quantity float64) error {
	if quantity < 0 {
		return errors.New(fmt.Sprintf("negative buy quantity %.8f", quantity))
	}

	if quantity > a.Assets[coin] {
		return errors.New(fmt.Sprintf("sell quantity: %.8f exceeds available: %.8f", quantity, a.Assets[coin]))
	}

	sellOrder := model.OrderRequest{
		Symbol:    coin,
		Side:      model.SELL,
		Type:      model.MARKET,
		Timestamp: a.GetTimeStamp(),
		Quantity:  quantity,
	}

	if err := a.Market.NewOrder(sellOrder); err != nil {
		return err
	}

	positionBuyValues := make([]int64, 0, len(a.Positions[coin]))
	for pbv := range a.Positions[coin] {
		positionBuyValues = append(positionBuyValues, pbv)
	}

	sort.Slice(positionBuyValues, func(i, j int) bool {
		return positionBuyValues[i] < positionBuyValues[j]
	})

	remainingQty := quantity

	for _, pbv := range positionBuyValues {
		if remainingQty >= a.Positions[coin][pbv] {
			remainingQty -= a.Positions[coin][pbv]
			delete(a.Positions[coin], pbv)
		} else {
			a.Positions[coin][pbv] -= remainingQty
			remainingQty = 0
		}
		if remainingQty == 0 {
			break
		}
	}

	a.Assets[coin] -= quantity

	currentValue := a.AssetValues[coin]
	a.Balance += model.IntFloatMul(currentValue, quantity*(1-a.Fee))

	return nil
}

func (a *Accountant) UpdateAssetValue(coin string, value int64, timestamp time.Time) error {
	if value < 0 {
		return errors.New(fmt.Sprintf("negative asset value (%s,%s)", coin, model.IntToString(value)))
	}
	a.AssetValues[coin] = value
	a.Market.UpdateCoinValue(coin, model.IntToFloat(value))
	return nil
}

func (a *Accountant) GetBalance() int64 {
	return a.Balance
}

func (a *Accountant) GetPositions(coin string) map[int64]float64 {
	return a.Positions[coin]
}

func (a *Accountant) GetFee() float64 {
	return a.Fee
}

func (a *Accountant) TotalAssetValue() int64 {
	var totalValue int64 = 0
	for coin, assetValue := range a.Assets {
		totalValue += model.IntFloatMul(a.AssetValues[coin], assetValue)
	}
	return totalValue
}

func (a *Accountant) NetWorth() int64 {
	return a.Balance + a.TotalAssetValue()
}

func (a *Accountant) AssetQty(asset string) float64 {
	if value, contains := a.Assets[asset]; contains {
		return value
	} else {
		return 0.0
	}
}

func (a *Accountant) AssetValue(asset string) int64 {
	return model.IntFloatMul(a.AssetValues[asset], a.AssetQty(asset))
}

func (a *Accountant) ToString() string {
	walletStr := fmt.Sprintf(">> NW:%s Balance:%s |", model.IntToString(a.NetWorth()), model.IntToString(a.GetBalance()))

	var coinList []string

	for coin, _ := range a.AssetValues {
		coinList = append(coinList, coin)
	}

	sort.Strings(coinList)

	for _, coin := range coinList {
		walletStr += fmt.Sprintf(" %s #%d Total:%s$(%.2f%%) |", coin, len(a.Positions[coin]),
			model.IntToString(a.AssetValue(coin)), float64(a.AssetValue(coin))/float64(a.NetWorth())*100)
	}

	return walletStr
}

func (a *Accountant) GetTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (a *Accountant) SyncWithMarket() {
	marketBalance, err := a.Market.Balance("USDT")

	if err == nil && a.Balance != model.FloatToInt(marketBalance.Free) {
		panic(fmt.Sprintf("Incoherent balance Acc:%d Market:%f", a.Balance, marketBalance.Free))
	}

}
