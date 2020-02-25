package model

import "github.com/shopspring/decimal"

type AccountInformation struct {
	MakerCommission  int64
	TakerCommission  int64
	BuyerCommission  int64
	SellerCommission int64
	CanTrade         bool
	CanWithdraw      bool
	CanDeposit       bool
	UpdateTime       int64
	AccountType      string
	Balances         []Balance
}

type Balance struct {
	Asset  string
	Free   decimal.Decimal
	Locked decimal.Decimal
}

type Trade struct {
	Symbol          string
	Id              int64
	OrderId         int64
	OrderListId     int64
	Price           decimal.Decimal
	Qty             decimal.Decimal
	Commission      decimal.Decimal
	CommissionAsset string
	Time            int64
	IsBuyer         bool
	IsMaker         bool
	IsBestMatch     bool
}
