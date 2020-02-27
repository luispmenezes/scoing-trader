package market

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"scoing-trader/trader/model/market/model"
	"sort"
	"time"
)

type Accountant struct {
	InitialBalance decimal.Decimal
	Fee            decimal.Decimal
	Balance        decimal.Decimal
	Market         model.Market
	Positions      map[string]map[string]decimal.Decimal
	Assets         map[string]decimal.Decimal
	AssetValues    map[string]decimal.Decimal
}

func NewAccountant(market model.Market, initialBalance decimal.Decimal, fee decimal.Decimal) *Accountant {
	return &Accountant{
		InitialBalance: initialBalance,
		Fee:            fee,
		Balance:        initialBalance,
		Market:         market,
		Positions:      make(map[string]map[string]decimal.Decimal),
		Assets:         make(map[string]decimal.Decimal),
		AssetValues:    make(map[string]decimal.Decimal),
	}
}

func (a *Accountant) Buy(coin string, quantity decimal.Decimal) (decimal.Decimal, error) {

	if quantity.LessThan(decimal.Zero) {
		return decimal.Zero, errors.New(fmt.Sprintf("negative buy quantity %s", quantity))
	}

	transactionValue := a.AssetValues[coin].Mul(quantity).Mul(a.Fee.Add(decimal.NewFromInt(1)))

	if transactionValue.GreaterThan(a.Balance) {
		return decimal.Zero, errors.New(fmt.Sprintf("buy transaction: %s exceeds balance: %s", transactionValue, a.Balance))
	}

	buyOrder := model.OrderRequest{
		Symbol:    coin,
		Side:      model.BUY,
		Type:      model.MARKET,
		Timestamp: a.GetTimeStamp(),
		Quantity:  quantity,
	}

	if err := a.Market.NewOrder(buyOrder); err != nil {
		return decimal.Zero, err
	}

	a.Balance = a.Balance.Sub(transactionValue)

	if _, hasKey := a.Positions[coin]; !hasKey {
		a.Positions[coin] = make(map[string]decimal.Decimal)
	}

	if _, hasKey := a.Positions[coin][a.AssetValues[coin].String()]; hasKey {
		a.Positions[coin][a.AssetValues[coin].String()] = a.Positions[coin][a.AssetValues[coin].String()].Add(quantity)
	} else {
		a.Positions[coin][a.AssetValues[coin].String()] = quantity
	}

	if _, hasKey := a.Assets[coin]; !hasKey {
		a.Assets[coin] = quantity
	} else {
		a.Assets[coin] = a.Assets[coin].Add(quantity)
	}

	return a.AssetValues[coin].Mul(quantity.Mul(decimal.NewFromInt(1).Add(a.Fee))), nil
}

func (a *Accountant) Sell(coin string, quantity decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	if quantity.LessThan(decimal.Zero) {
		return decimal.Zero, decimal.Zero, errors.New(fmt.Sprintf("negative buy quantity %s", quantity))
	}

	if quantity.GreaterThan(a.Assets[coin]) {
		return decimal.Zero, decimal.Zero, errors.New(fmt.Sprintf("sell quantity: %s exceeds available: %s", quantity, a.Assets[coin]))
	}

	sellOrder := model.OrderRequest{
		Symbol:    coin,
		Side:      model.SELL,
		Type:      model.MARKET,
		Timestamp: a.GetTimeStamp(),
		Quantity:  quantity,
	}

	if err := a.Market.NewOrder(sellOrder); err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	positionBuyValues := make([]string, 0, len(a.Positions[coin]))
	for pbv := range a.Positions[coin] {
		positionBuyValues = append(positionBuyValues, pbv)
	}

	sort.Slice(positionBuyValues, func(i, j int) bool {
		decimalI, _ := decimal.NewFromString(positionBuyValues[i])
		decimalJ, _ := decimal.NewFromString(positionBuyValues[j])

		return decimalI.LessThan(decimalJ)
	})

	remainingQty := quantity
	var positionTransactionSum decimal.Decimal

	for _, pbv := range positionBuyValues {
		if remainingQty.GreaterThanOrEqual(a.Positions[coin][pbv]) {
			remainingQty = remainingQty.Sub(a.Positions[coin][pbv])
			pbvDecimal, _ := decimal.NewFromString(pbv)
			positionTransactionSum = positionTransactionSum.Add(a.Positions[coin][pbv].Mul(pbvDecimal).Mul(decimal.NewFromInt(1).Add(a.Fee)))
			delete(a.Positions[coin], pbv)
		} else {
			a.Positions[coin][pbv] = a.Positions[coin][pbv].Sub(remainingQty)
			remainingQty = decimal.Zero
		}
		if remainingQty.IsZero() {
			break
		}
	}

	a.Assets[coin] = a.Assets[coin].Sub(quantity)
	transaction := a.AssetValues[coin].Mul(quantity.Mul(decimal.NewFromInt(1).Sub(a.Fee)))
	profit := transaction.Sub(positionTransactionSum)
	a.Balance = a.Balance.Add(transaction)

	return transaction, profit, nil
}

func (a *Accountant) UpdateAssetValue(coin string, value decimal.Decimal, timestamp time.Time) error {
	if value.LessThan(decimal.Zero) {
		return errors.New(fmt.Sprintf("negative asset value (%s,%s)", coin, value))
	}
	a.AssetValues[coin] = value
	a.Market.UpdateCoinValue(coin, value)
	return nil
}

func (a *Accountant) GetBalance() decimal.Decimal {
	return a.Balance
}

func (a *Accountant) GetPositions(coin string) map[string]decimal.Decimal {
	return a.Positions[coin]
}

func (a *Accountant) GetFee() decimal.Decimal {
	return a.Fee
}

func (a *Accountant) TotalAssetValue() decimal.Decimal {
	var totalValue decimal.Decimal
	for coin, assetValue := range a.Assets {
		totalValue = totalValue.Add(a.AssetValues[coin].Mul(assetValue))
	}
	return totalValue
}

func (a *Accountant) NetWorth() decimal.Decimal {
	return a.Balance.Add(a.TotalAssetValue())
}

func (a *Accountant) AssetQty(asset string) decimal.Decimal {
	if value, contains := a.Assets[asset]; contains {
		return value
	} else {
		return decimal.Zero
	}
}

func (a *Accountant) AssetValue(asset string) decimal.Decimal {
	return a.AssetValues[asset].Mul(a.AssetQty(asset))
}

func (a *Accountant) ToString() string {
	nw, _ := a.NetWorth().Float64()
	balance, _ := a.GetBalance().Float64()

	walletStr := fmt.Sprintf(">> NW:%.4f Balance:%.4f |", nw, balance)

	var coinList []string

	for coin, _ := range a.AssetValues {
		coinList = append(coinList, coin)
	}

	sort.Strings(coinList)

	for _, coin := range coinList {
		assetValue, _ := a.AssetValue(coin).Float64()
		assetPercentage, _ := a.AssetValue(coin).Div(a.NetWorth().Mul(decimal.NewFromInt(100))).Float64()
		walletStr += fmt.Sprintf(" %s #%d Total:%.4f$(%4.f%%) |", coin, len(a.Positions[coin]),
			assetValue, assetPercentage)
	}

	return walletStr
}

func (a *Accountant) GetTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (a *Accountant) SyncWithMarket() {
	marketBalance, err := a.Market.Balance("USDT")

	if err == nil && !a.Balance.Equal(marketBalance.Free) {
		panic(fmt.Sprintf("Incoherent balance Acc:%s Market:%s", a.Balance, marketBalance.Free))
	}
}
