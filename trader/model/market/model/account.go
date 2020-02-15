package model

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
	Free   float64
	Locked float64
}

type Trade struct {
	Symbol          string
	Id              int64
	OrderId         int64
	OrderListId     int64
	Price           float64
	Qty             float64
	Commission      float64
	CommissionAsset string
	Time            int64
	IsBuyer         bool
	IsMaker         bool
	IsBestMatch     bool
}
