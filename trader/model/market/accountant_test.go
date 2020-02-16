package market

import (
	"fmt"
	"testing"
	"time"
)

func TestBuy(t *testing.T) {
	market := NewSimulatedMarket(0, 0)
	market.Deposit("USDT", 100)
	accountant := NewAccountant(market, 100, 0)
	err := accountant.UpdateCoinValue("BTCUSDT", 10, time.Now())

	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", 2)

	if err != nil {
		t.Error(err)
	}

	if accountant.Balance != 80 {
		t.Error("Expected balance=80, got ", accountant.Balance)
	}

	positionsByCoin, hasKey := accountant.Positions["BTCUSDT"]
	if !hasKey {
		t.Error("Missing coin in positions")
	}

	if positionsByCoin[10] != 2 {
		t.Error("Expected position=2, got ", positionsByCoin[2])
	}
}

func TestSell(t *testing.T) {
	market := NewSimulatedMarket(0, 0)
	market.Deposit("USDT", 100)
	accountant := NewAccountant(market, 100, 0)
	err := accountant.UpdateCoinValue("BTCUSDT", 10, time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", 2)
	if err != nil {
		t.Error(err)
	}

	err = accountant.UpdateCoinValue("BTCUSDT", 20, time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", 2)
	if err != nil {
		t.Error(err)
	}

	if accountant.Balance != 40 {
		t.Error("Expected balance before sale 40, got ", accountant.Balance)
	}

	if accountant.CoinPortfolioValue("BTCUSDT") != 4 {
		t.Error("Expected BTC balance 4, got ", accountant.CoinPortfolioValue("BTCUSDT"))
	}

	err = accountant.Sell("BTCUSDT", 3)

	if err != nil {
		t.Error(err)
	}

	if len(accountant.Positions["BTCUSDT"]) != 1 || accountant.Positions["BTCUSDT"][20] != 1 {
		t.Error(fmt.Sprintf("Invalid positions length:(expected 1 got %d) qty@20(expected: 1 got %f)",
			len(accountant.Positions["BTCUSDT"]), accountant.Positions["BTCUSDT"][20]))
	}

	err = accountant.Sell("BTCUSDT", 1)
	if err != nil {
		t.Error(err)
	}

	if accountant.Balance != 120 {
		t.Error("Expected balance=120, got ", accountant.Balance)
	}

	if len(accountant.Positions["BTCUSDT"]) != 0 {
		t.Error("Positions still open")
	}
}

func TestFee(t *testing.T) {
	market := NewSimulatedMarket(0, 0.1)
	market.Deposit("USDT", 100)
	accountant := NewAccountant(market, 100, 0.1)
	err := accountant.UpdateCoinValue("BTCUSDT", 10, time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", 1)
	if err != nil {
		t.Error(err)
	}

	err = accountant.Sell("BTCUSDT", 1)
	if err != nil {
		t.Error(err)
	}

	if accountant.Balance != 98 {
		t.Error("Expected balance=98, got ", accountant.Balance)
	}
}

func TestPositionValue(t *testing.T) {
	market := NewSimulatedMarket(0, 0)
	market.Deposit("USDT", 100)
	accountant := NewAccountant(market, 100, 0)
	err := accountant.UpdateCoinValue("BTCUSDT", 10, time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.UpdateCoinValue("ETHUSDT", 5, time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", 2)
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("ETHUSDT", 10)
	if err != nil {
		t.Error(err)
	}

	if accountant.NetWorth() != 100 {
		t.Error("Expected Position Value 100, got ", accountant.NetWorth())
	}

	err = accountant.UpdateCoinValue("BTCUSDT", 20, time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.UpdateCoinValue("ETHUSDT", 10, time.Now())
	if err != nil {
		t.Error(err)
	}

	if accountant.NetWorth() != 170 {
		t.Error("Expected Position Value 170, got ", accountant.NetWorth())
	}
}
