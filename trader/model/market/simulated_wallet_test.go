package market

import (
	"testing"
	"time"
)

func TestBuy(t *testing.T) {
	wallet := NewSimulatedWallet(100, 0)
	wallet.UpdateCoinValue("BTC", 10, time.Now())
	wallet.Buy("BTC", 2)
	if wallet.Balance != 80 {
		t.Error("Expected balance=80, got ", wallet.Balance)
	}
	positionsByCoin, hasKey := wallet.Positions["BTC"]
	if !hasKey {
		t.Error("Missing coin in positions")
	}
	if positionsByCoin[10] != 2 {
		t.Error("Expected position=2, got ", positionsByCoin[2])
	}
}

func TestSell(t *testing.T) {
	wallet := NewSimulatedWallet(100, 0)
	wallet.UpdateCoinValue("BTC", 10, time.Now())
	wallet.Buy("BTC", 2)
	wallet.UpdateCoinValue("BTC", 20, time.Now())
	wallet.Buy("BTC", 2)

	if wallet.Balance != 40 {
		t.Error("Expected balance before sale 40, got ", wallet.Balance)
	}

	wallet.Sell("BTC", 10, 2)
	wallet.Sell("BTC", 20, 2)

	if wallet.Balance != 120 {
		t.Error("Expected balance=120, got ", wallet.Balance)
	}

	if len(wallet.Positions["BTC"]) != 0 {
		t.Error("Positions still open")
	}
}

func TestFee(t *testing.T) {
	wallet := NewSimulatedWallet(100, 0.1)
	wallet.UpdateCoinValue("BTC", 10, time.Now())
	wallet.Buy("BTC", 1)
	wallet.Sell("BTC", 10, 1)
	if wallet.Balance != 98 {
		t.Error("Expected balance=98, got ", wallet.Balance)
	}
}

func TestPositionValue(t *testing.T) {
	wallet := NewSimulatedWallet(100, 0)
	wallet.UpdateCoinValue("BTC", 10, time.Now())
	wallet.UpdateCoinValue("ETH", 5, time.Now())
	wallet.Buy("BTC", 2)
	wallet.Buy("ETH", 10)
	wallet.UpdateCoinValue("BTC", 20, time.Now())
	wallet.UpdateCoinValue("ETH", 10, time.Now())
	if wallet.NetWorth() != 170 {
		t.Error("Expected Position Value 170, got ", wallet.NetWorth())
	}
}
