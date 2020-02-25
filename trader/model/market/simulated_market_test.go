package market

import (
	"fmt"
	"github.com/shopspring/decimal"
	"scoing-trader/trader/model/market/model"
	"testing"
)

func TestMarketDeposit(t *testing.T) {
	market := NewSimulatedMarket(0, decimal.NewFromFloat(0.001))
	market.Deposit("USDT", decimal.NewFromInt(1000))

	balance, err := market.Balance("USDT")

	if err != nil {
		t.Error("Coin balance not found")
	} else {
		if !balance.Free.Equal(decimal.NewFromInt(1000)) {
			t.Error(fmt.Sprintf("Incorrect free balance got %s exepected %f", balance.Free, 1000.0))
		}
		if !balance.Locked.IsZero() {
			t.Error(fmt.Sprintf("Incorrect locked balance got %s exepected %f", balance.Locked, 0.0))
		}
	}
}

func TestMarketBuy(t *testing.T) {
	market := NewSimulatedMarket(0, decimal.NewFromFloat(0.001))
	market.Deposit("USDT", decimal.NewFromInt(1000))
	order := model.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          model.BUY,
		Type:          model.MARKET,
		Timestamp:     0,
		TimeInForce:   model.GTC,
		Quantity:      decimal.NewFromFloat(0.5),
		QuoteOrderQty: decimal.NewFromInt(500),
		Price:         decimal.NewFromInt(1000),
		ClientOrderId: "1234",
		StopPrice:     decimal.NewFromInt(1005),
		IcebergQty:    decimal.Zero,
		ResponseType:  model.FULL,
	}

	err := market.NewOrder(order)

	if err != nil {
		t.Error(err)
	}

	if len(market.OrderHistory()) != 1 {
		t.Error(fmt.Sprintf("Inconsistent order history expected: 1 got: %d", len(market.OrderHistory())))
	}

	if len(market.Trades()) != 1 {
		t.Error(fmt.Sprintf("Inconsistent trade history expected: 1 got: %d", len(market.Trades())))
	}

	if len(market.OpenOrders("BTCUSDT")) > 0 {
		t.Error("Order is open")
	}

	balanceBTC, errBTC := market.Balance("BTC")
	if errBTC != nil {
		t.Error("Coin balance not found")
	} else {
		if !balanceBTC.Free.Equal(decimal.NewFromFloat(0.5)) {
			t.Error(fmt.Sprintf("Incorrect coin free balance got %s exepected %f", balanceBTC.Free, 0.5))
		}
	}

	balanceUSDT, errUSDT := market.Balance("USDT")
	if errUSDT != nil {
		t.Error("Coin balance not found")
	} else {
		if !balanceUSDT.Free.Equal(decimal.NewFromFloat(499.5)) {
			t.Error(fmt.Sprintf("Incorrect coin free balance got %s exepected %f", balanceUSDT.Free, 499.5))
		}
	}
}

func TestMarketUnfilled(t *testing.T) {
	market := NewSimulatedMarket(1, decimal.NewFromFloat(0.001))
	market.Deposit("USDT", decimal.NewFromInt(1000))
	order1 := model.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          model.BUY,
		Type:          model.MARKET,
		Timestamp:     0,
		TimeInForce:   model.GTC,
		Quantity:      decimal.NewFromFloat(0.1),
		QuoteOrderQty: decimal.NewFromInt(100),
		Price:         decimal.NewFromInt(1000),
		ClientOrderId: "1",
		StopPrice:     decimal.NewFromInt(1005),
		IcebergQty:    decimal.Zero,
		ResponseType:  model.FULL,
	}

	err := market.NewOrder(order1)

	if err != nil {
		t.Error("Order 1 failed")
	}

	if len(market.OpenOrders("BTCUSDT")) < 0 {
		t.Error("Order is not open ")
	}

	order2 := model.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          model.BUY,
		Type:          model.MARKET,
		Timestamp:     0,
		TimeInForce:   model.GTC,
		Quantity:      decimal.NewFromFloat(0.1),
		QuoteOrderQty: decimal.NewFromInt(100),
		Price:         decimal.NewFromInt(1000),
		ClientOrderId: "2",
		StopPrice:     decimal.NewFromInt(1005),
		IcebergQty:    decimal.Zero,
		ResponseType:  model.FULL,
	}

	err2 := market.NewOrder(order2)

	if err2 == nil {
		t.Error("Order 2 was accepted")
	}

	btcBalance, errBalBTC := market.Balance("BTC")

	if errBalBTC != nil {
		t.Error(errBalBTC)
	} else if !btcBalance.Free.IsZero() {
		t.Error(fmt.Sprintf("Inavlid BTC balance. Expected 0 got %s", btcBalance.Free))
	}

	usdtBalance, errBalUSDT := market.Balance("USDT")

	if errBalUSDT != nil {
		t.Error(usdtBalance)
	} else {
		if !usdtBalance.Free.Equal(decimal.NewFromFloat(899.9)) {
			t.Error(fmt.Sprintf("Inavlid USDT balance. Expected 899.90 got %s", btcBalance.Free))
		}
		if !usdtBalance.Locked.Equal(decimal.NewFromFloat(100.1)) {
			t.Error(fmt.Sprintf("Incorrect locked balance. Expected 100.10 got %s", usdtBalance.Locked))
		}
	}

	errCancel := market.CancelOrder("1")

	if errCancel != nil {
		t.Error(errCancel)
	} else {
		usdtBalance, _ := market.Balance("USDT")

		if !usdtBalance.Free.Equal(decimal.NewFromInt(1000)) {
			t.Error(fmt.Sprintf("Inavlid USDT balance. Expected 1000.00 got %s", usdtBalance.Free))
		}
		if !usdtBalance.Locked.IsZero() {
			t.Error(fmt.Sprintf("Incorrect locked balance. Expected 0.00 got %s", usdtBalance.Locked))
		}
	}
}
