package market

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestBuy(t *testing.T) {
	market := NewSimulatedMarket(0, decimal.Zero)
	market.Deposit("USDT", decimal.NewFromInt(100))
	accountant := NewAccountant(market, decimal.NewFromInt(100), decimal.Zero)
	err := accountant.UpdateAssetValue("BTCUSDT", decimal.NewFromInt(10), time.Now())

	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", decimal.NewFromInt(2))

	if err != nil {
		t.Error(err)
	}

	if !accountant.Balance.Equal(decimal.NewFromInt(80)) {
		t.Error("Expected balance=80, got ", accountant.Balance)
	}

	positionsByCoin, hasKey := accountant.Positions["BTCUSDT"]
	if !hasKey {
		t.Error("Missing coin in positions")
	}

	if !positionsByCoin[decimal.NewFromInt(10).String()].Equal(decimal.NewFromInt(2)) {
		t.Error("Expected position=2, got ", positionsByCoin[decimal.NewFromInt(2).String()])
	}
}

func TestSell(t *testing.T) {
	market := NewSimulatedMarket(0, decimal.Zero)
	market.Deposit("USDT", decimal.NewFromInt(100))
	accountant := NewAccountant(market, decimal.NewFromInt(100), decimal.Zero)
	err := accountant.UpdateAssetValue("BTCUSDT", decimal.NewFromInt(10), time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", decimal.NewFromInt(2))
	if err != nil {
		t.Error(err)
	}

	err = accountant.UpdateAssetValue("BTCUSDT", decimal.NewFromInt(20), time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", decimal.NewFromInt(2))
	if err != nil {
		t.Error(err)
	}

	if !accountant.Balance.Equal(decimal.NewFromInt(40)) {
		t.Error("Expected balance before sale 40, got ", accountant.Balance)
	}

	if !accountant.AssetQty("BTCUSDT").Equal(decimal.NewFromInt(4)) {
		t.Error("Expected BTC balance 4, got ", accountant.AssetQty("BTCUSDT"))
	}

	err = accountant.Sell("BTCUSDT", decimal.NewFromInt(3))

	if err != nil {
		t.Error(err)
	}

	if len(accountant.Positions["BTCUSDT"]) != 1 || !accountant.Positions["BTCUSDT"][decimal.NewFromInt(20).String()].Equal(decimal.NewFromInt(1)) {
		t.Error(fmt.Sprintf("Invalid positions length:(expected 1 got %d) qty@20(expected: 1 got %s)",
			len(accountant.Positions["BTCUSDT"]), accountant.Positions["BTCUSDT"][decimal.NewFromInt(20).String()]))
	}

	err = accountant.Sell("BTCUSDT", decimal.NewFromInt(1))
	if err != nil {
		t.Error(err)
	}

	if !accountant.Balance.Equal(decimal.NewFromInt(120)) {
		t.Error("Expected balance=120, got ", accountant.Balance)
	}

	if len(accountant.Positions["BTCUSDT"]) != 0 {
		t.Error("Positions still open")
	}
}

func TestFee(t *testing.T) {
	market := NewSimulatedMarket(0, decimal.NewFromFloat(0.1))
	market.Deposit("USDT", decimal.NewFromInt(100))
	accountant := NewAccountant(market, decimal.NewFromInt(100), decimal.NewFromFloat(0.1))
	err := accountant.UpdateAssetValue("BTCUSDT", decimal.NewFromInt(10), time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", decimal.NewFromInt(1))
	if err != nil {
		t.Error(err)
	}

	err = accountant.Sell("BTCUSDT", decimal.NewFromInt(1))
	if err != nil {
		t.Error(err)
	}

	if !accountant.Balance.Equal(decimal.NewFromInt(98)) {
		t.Error("Expected balance=98, got ", accountant.Balance)
	}
}

func TestPositionValue(t *testing.T) {
	market := NewSimulatedMarket(0, decimal.Zero)
	market.Deposit("USDT", decimal.NewFromInt(100))
	accountant := NewAccountant(market, decimal.NewFromInt(100), decimal.Zero)
	err := accountant.UpdateAssetValue("BTCUSDT", decimal.NewFromInt(10), time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.UpdateAssetValue("ETHUSDT", decimal.NewFromInt(5), time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("BTCUSDT", decimal.NewFromInt(2))
	if err != nil {
		t.Error(err)
	}

	err = accountant.Buy("ETHUSDT", decimal.NewFromInt(10))
	if err != nil {
		t.Error(err)
	}

	if !accountant.NetWorth().Equal(decimal.NewFromInt(100)) {
		t.Error("Expected Position Value 100, got ", accountant.NetWorth())
	}

	err = accountant.UpdateAssetValue("BTCUSDT", decimal.NewFromInt(20), time.Now())
	if err != nil {
		t.Error(err)
	}

	err = accountant.UpdateAssetValue("ETHUSDT", decimal.NewFromInt(10), time.Now())
	if err != nil {
		t.Error(err)
	}

	if !accountant.NetWorth().Equal(decimal.NewFromInt(170)) {
		t.Error("Expected Position Value 170, got ", accountant.NetWorth())
	}
}
