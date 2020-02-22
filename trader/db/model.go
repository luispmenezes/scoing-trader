package db

import "time"

type TraderEvent struct {
	tableName struct{}  `sql:"cointron.trader_events"`
	Id        string    `sql:"event_id,pk"`
	EventTime time.Time `sql:"event_time"`
	Type      string    `sql:"event_type"`
	Coin      string    `sql:"coin"`
	Balance   int64     `sql:"balance,use_zero"`
	Networth  int64     `sql:"net_worth,use_zero"`
	Positions string    `sql:"positions"`
}

type TraderConfig struct {
	tableName struct{}  `sql:"cointron.trader_config"`
	Id        string    `sql:"config_id,pk"`
	EventTime time.Time `sql:"creation_time"`
	Strategy  string    `sql:"strategy_name"`
	Fitness   int64     `sql:"fitness,use_zero"`
}
