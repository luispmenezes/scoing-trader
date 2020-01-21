package wallet

import (
	"fmt"
	"sort"
	"time"
)

type SimulatedWallet struct {
	InitialBalance float64
	Fee            float64
	Balance        float64
	Positions      map[string]map[float64]float64
	CoinValues     map[string]float64
}

func NewSimulatedWallet(initialBalance float64, fee float64) *SimulatedWallet {
	return &SimulatedWallet{
		InitialBalance: initialBalance,
		Fee:            fee,
		Balance:        initialBalance,
		Positions:      make(map[string]map[float64]float64),
		CoinValues:     make(map[string]float64),
	}
}

func (w *SimulatedWallet) Buy(coin string, quantity float64) {
	if quantity > 0 {
		currentValue := w.CoinValues[coin]
		transaction := currentValue * quantity
		transactionFee := transaction * w.Fee
		totalCost := transaction + transactionFee
		if w.Balance > totalCost {
			w.Balance -= totalCost
			_, hasKey := w.Positions[coin]
			if !hasKey {
				w.Positions[coin] = make(map[float64]float64)
			}
			_, hasKey = w.Positions[coin][currentValue]
			if hasKey {
				w.Positions[coin][currentValue] += quantity
			} else {
				w.Positions[coin][currentValue] = quantity
			}
		}
	}
}

func (w *SimulatedWallet) Sell(coin string, buyValue float64, quantity float64) {
	if quantity > 0 && quantity <= w.Positions[coin][buyValue] {
		currentValue := w.CoinValues[coin]
		w.Balance += currentValue * quantity * (1 - w.Fee)
		delete(w.Positions[coin], buyValue)
		if len(w.Positions[coin]) == 0 {
			delete(w.Positions, coin)
		}
	}
}

func (w *SimulatedWallet) UpdateCoinValue(coin string, value float64, timestamp time.Time) {
	if value >= 0 {
		w.CoinValues[coin] = value
	} else {
		panic("Negative coin value")
	}
}

func (w *SimulatedWallet) GetBalance() float64 {
	return w.Balance
}

func (w *SimulatedWallet) GetPositions(coin string) map[float64]float64 {
	return w.Positions[coin]
}

func (w *SimulatedWallet) GetFee() float64 {
	return w.Fee
}

func (w *SimulatedWallet) TotalPositionValue() float64 {
	var totalValue float64 = 0
	for coin, coinValue := range w.CoinValues {
		_, hasKey := w.Positions[coin]
		if hasKey {
			for _, positionQty := range w.Positions[coin] {
				totalValue += positionQty * coinValue
			}
		}
	}
	return totalValue
}

func (w *SimulatedWallet) NetWorth() float64 {
	return w.Balance + w.TotalPositionValue()
}

func (w *SimulatedWallet) CoinNetWorth(coin string) float64 {
	totalQty := 0.0
	for _, qty := range w.GetPositions(coin) {
		totalQty += qty
	}
	return totalQty * w.CoinValues[coin]
}

func (w *SimulatedWallet) ToString() string {
	walletStr := fmt.Sprintf(">> NW:%.2f Balance:%.2f |", w.NetWorth(), w.GetBalance())

	var coinList []string

	for coin, _ := range w.CoinValues {
		coinList = append(coinList, coin)
	}

	sort.Strings(coinList)

	for _, coin := range coinList {
		walletStr += fmt.Sprintf(" %s #%d Total:%.2f(%.2f%%) |", coin, len(w.Positions[coin]), w.CoinNetWorth(coin),
			w.CoinNetWorth(coin)/w.NetWorth()*100)
	}

	return walletStr
}
