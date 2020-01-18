package trader

import (
	"math/rand"
)

type TraderConfig struct {
	BuyPred15Mod   float64
	BuyPred60Mod   float64
	BuyPred1440Mod float64
	StopLoss       float64
	ProfitCap      float64
	BuyNWQtyMod    float64
	BuyQty15Mod    float64
	BuyQty60Mod    float64
	BuyQty1440Mod  float64
	SellPosQtyMod  float64
	SellQty15Mod   float64
	SellQty60Mod   float64
	SellQty1440Mod float64
}

func (t *TraderConfig) NumParams() int {
	return 12
}

func RandomConfig() *TraderConfig {
	return &TraderConfig{
		BuyPred15Mod:   randomFloat(-1, 1),
		BuyPred60Mod:   randomFloat(-1, 1),
		BuyPred1440Mod: randomFloat(-1, 1),
		StopLoss:       randomFloat(-0.2, 0),
		ProfitCap:      randomFloat(0, 0.1),
		BuyNWQtyMod:    randomFloat(-1, 1),
		BuyQty15Mod:    randomFloat(-1, 1),
		BuyQty60Mod:    randomFloat(-1, 1),
		BuyQty1440Mod:  randomFloat(-1, 1),
		SellPosQtyMod:  randomFloat(-1, 1),
		SellQty15Mod:   randomFloat(-1, 1),
		SellQty60Mod:   randomFloat(-1, 1),
		SellQty1440Mod: randomFloat(-1, 1),
	}
}

func RandomBetweenTwo(a, b TraderConfig) *TraderConfig {
	return &TraderConfig{
		BuyPred15Mod:   randomFloat(a.BuyPred15Mod, b.BuyPred15Mod),
		BuyPred60Mod:   randomFloat(a.BuyPred60Mod, b.BuyPred60Mod),
		BuyPred1440Mod: randomFloat(a.BuyPred1440Mod, b.BuyPred1440Mod),
		StopLoss:       randomFloat(a.StopLoss, b.StopLoss),
		ProfitCap:      randomFloat(a.ProfitCap, b.ProfitCap),
		BuyNWQtyMod:    randomFloat(a.BuyNWQtyMod, b.BuyNWQtyMod),
		BuyQty15Mod:    randomFloat(a.BuyQty15Mod, b.BuyQty15Mod),
		BuyQty60Mod:    randomFloat(a.BuyQty60Mod, b.BuyQty60Mod),
		BuyQty1440Mod:  randomFloat(a.BuyQty1440Mod, b.BuyQty1440Mod),
		SellPosQtyMod:  randomFloat(a.SellPosQtyMod, b.SellPosQtyMod),
		SellQty15Mod:   randomFloat(a.SellQty15Mod, b.SellQty15Mod),
		SellQty60Mod:   randomFloat(a.SellQty60Mod, b.SellQty60Mod),
		SellQty1440Mod: randomFloat(a.SellQty1440Mod, b.SellQty1440Mod),
	}
}

func (t *TraderConfig) RandomizeParam() *TraderConfig {
	idx := randomFloat(0, float64(t.NumParams()))
	switch idx {
	case 0:
		t.BuyPred15Mod = randomFloat(-1, 1)
		break
	case 1:
		t.BuyPred60Mod = randomFloat(-1, 1)
		break
	case 2:
		t.BuyPred1440Mod = randomFloat(-1, 1)
		break
	case 3:
		t.StopLoss = randomFloat(-0.3, 0)
		break
	case 4:
		t.ProfitCap = randomFloat(0, 0.3)
		break
	case 5:
		t.BuyNWQtyMod = randomFloat(-1, 1)
		break
	case 6:
		t.BuyQty15Mod = randomFloat(-1, 1)
		break
	case 7:
		t.BuyQty60Mod = randomFloat(-1, 1)
		break
	case 8:
		t.BuyQty1440Mod = randomFloat(-1, 1)
		break
	case 9:
		t.SellPosQtyMod = randomFloat(-1, 1)
		break
	case 10:
		t.SellQty15Mod = randomFloat(-1, 1)
		break
	case 11:
		t.SellQty60Mod = randomFloat(-1, 1)
		break
	case 12:
		t.SellQty1440Mod = randomFloat(-1, 1)
		break
	}
	return t

}

func randomFloat(a, b float64) float64 {
	var min, max float64

	if a > b {
		max = a
		min = b
	} else {
		max = b
		min = a
	}

	return min + rand.Float64()*(max-min)
}
