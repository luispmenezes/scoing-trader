package market

import (
	"fmt"
	"scoing-trader/trader/model/market/model"
	"testing"
)

func TestMarketDeposit(t *testing.T) {
	market := NewSimulatedMarket(0, 0.001)
	market.Deposit("USDT", 1000)

	balance, err := market.Balance("USDT")

	if err != nil {
		t.Error("Coin balance not found")
	} else {
		if balance.Free != 1000 {
			t.Error(fmt.Sprintf("Incorrect free balance got %f exepected %f", balance.Free, 1000.0))
		}
		if balance.Locked != 0 {
			t.Error(fmt.Sprintf("Incorrect locked balance got %f exepected %f", balance.Locked, 0.0))
		}
	}
}

func TestMarketBuy(t *testing.T) {
	market := NewSimulatedMarket(0, 0.001)
	market.Deposit("USDT", 1000)
	order := model.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          model.BUY,
		Type:          model.MARKET,
		Timestamp:     0,
		TimeInForce:   model.GTC,
		Quantity:      0.5,
		QuoteOrderQty: 500,
		Price:         1000,
		ClientOrderId: "1234",
		StopPrice:     1005,
		IcebergQty:    0,
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
		if balanceBTC.Free != 0.5 {
			t.Error(fmt.Sprintf("Incorrect coin free balance got %f exepected %f", balanceBTC.Free, 0.5))
		}
	}

	balanceUSDT, errUSDT := market.Balance("USDT")
	if errUSDT != nil {
		t.Error("Coin balance not found")
	} else {
		if fmt.Sprintf("%.4f", balanceUSDT.Free) != "499.5000" {
			t.Error(fmt.Sprintf("Incorrect coin free balance got %f exepected %f", balanceUSDT.Free, 499.5))
		}
	}
}

func TestMarketUnfilled(t *testing.T) {
	market := NewSimulatedMarket(1, 0.001)
	market.Deposit("USDT", 1000)
	order1 := model.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          model.BUY,
		Type:          model.MARKET,
		Timestamp:     0,
		TimeInForce:   model.GTC,
		Quantity:      0.1,
		QuoteOrderQty: 100,
		Price:         1000,
		ClientOrderId: "1",
		StopPrice:     1005,
		IcebergQty:    0,
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
		Quantity:      0.1,
		QuoteOrderQty: 100,
		Price:         1000,
		ClientOrderId: "2",
		StopPrice:     1005,
		IcebergQty:    0,
		ResponseType:  model.FULL,
	}

	err2 := market.NewOrder(order2)

	if err2 == nil {
		t.Error("Order 2 was accepted")
	}

	btcBalance, errBalBTC := market.Balance("BTC")

	if errBalBTC != nil {
		t.Error(errBalBTC)
	} else if btcBalance.Free != 0 {
		t.Error(fmt.Sprintf("Inavlid BTC balance. Expected 0 got %f", btcBalance.Free))
	}

	usdtBalance, errBalUSDT := market.Balance("USDT")

	if errBalUSDT != nil {
		t.Error(usdtBalance)
	} else {
		if fmt.Sprintf("%.2f", usdtBalance.Free) != "899.90" {
			t.Error(fmt.Sprintf("Inavlid USDT balance. Expected 899.90 got %f", btcBalance.Free))
		}
		if fmt.Sprintf("%.2f", usdtBalance.Locked) != "100.10" {
			t.Error(fmt.Sprintf("Incorrect locked balance. Expected 100.10 got %f", usdtBalance.Locked))
		}
	}

	errCancel := market.CancelOrder("1")

	if errCancel != nil {
		t.Error(errCancel)
	} else {
		usdtBalance, _ := market.Balance("USDT")

		if fmt.Sprintf("%.2f", usdtBalance.Free) != "1000.00" {
			t.Error(fmt.Sprintf("Inavlid USDT balance. Expected 1000.00 got %f", btcBalance.Free))
		}
		if fmt.Sprintf("%.2f", usdtBalance.Locked) != "0.00" {
			t.Error(fmt.Sprintf("Incorrect locked balance. Expected 0.00 got %f", usdtBalance.Locked))
		}
	}
}
