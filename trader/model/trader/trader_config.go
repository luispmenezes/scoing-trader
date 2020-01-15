package trader

import (
	"math/rand"
)

type TraderConfig struct {
	BuyPred15Mod    float64
	BuyPred60Mod    float64
	BuyPred1440Mod  float64
	SellPred15Mod   float64
	SellPred60Mod   float64
	SellPred1440Mod float64
	StopLoss        float64
	ProfitCap       float64
	BuyNWQtyMod     float64
	BuyQty15Mod     float64
	BuyQty60Mod     float64
	BuyQty1440Mod   float64
	SellPosQtyMod   float64
	SellQty15Mod    float64
	SellQty60Mod    float64
	SellQty1440Mod  float64
}

func (t *TraderConfig) NumParams() int {
	return 15
}

func RandomConfig(min, max float64) *TraderConfig {
	return &TraderConfig{
		BuyPred15Mod:    randomFloat(min, max),
		BuyPred60Mod:    randomFloat(min, max),
		BuyPred1440Mod:  randomFloat(min, max),
		SellPred15Mod:   randomFloat(min, max),
		SellPred60Mod:   randomFloat(min, max),
		SellPred1440Mod: randomFloat(min, max),
		StopLoss:        randomFloat(min, max),
		ProfitCap:       randomFloat(min, max),
		BuyNWQtyMod:     randomFloat(min, max),
		BuyQty15Mod:     randomFloat(min, max),
		BuyQty60Mod:     randomFloat(min, max),
		BuyQty1440Mod:   randomFloat(min, max),
		SellPosQtyMod:   randomFloat(min, max),
		SellQty15Mod:    randomFloat(min, max),
		SellQty60Mod:    randomFloat(min, max),
		SellQty1440Mod:  randomFloat(min, max),
	}
}

func RandomBetweenTwo(a, b TraderConfig) *TraderConfig {
	return &TraderConfig{
		BuyPred15Mod:    randomFloat(a.BuyPred15Mod, b.BuyPred15Mod),
		BuyPred60Mod:    randomFloat(a.BuyPred60Mod, b.BuyPred60Mod),
		BuyPred1440Mod:  randomFloat(a.BuyPred1440Mod, b.BuyPred1440Mod),
		SellPred15Mod:   randomFloat(a.SellPred15Mod, b.SellPred15Mod),
		SellPred60Mod:   randomFloat(a.SellPred60Mod, b.SellPred60Mod),
		SellPred1440Mod: randomFloat(a.SellPred1440Mod, b.SellPred1440Mod),
		StopLoss:        randomFloat(a.StopLoss, b.StopLoss),
		ProfitCap:       randomFloat(a.ProfitCap, b.ProfitCap),
		BuyNWQtyMod:     randomFloat(a.BuyNWQtyMod, b.BuyNWQtyMod),
		BuyQty15Mod:     randomFloat(a.BuyQty15Mod, b.BuyQty15Mod),
		BuyQty60Mod:     randomFloat(a.BuyQty60Mod, b.BuyQty60Mod),
		BuyQty1440Mod:   randomFloat(a.BuyQty1440Mod, b.BuyQty1440Mod),
		SellPosQtyMod:   randomFloat(a.SellPosQtyMod, b.SellPosQtyMod),
		SellQty15Mod:    randomFloat(a.SellQty15Mod, b.SellQty15Mod),
		SellQty60Mod:    randomFloat(a.SellQty60Mod, b.SellQty60Mod),
		SellQty1440Mod:  randomFloat(a.SellQty1440Mod, b.SellQty1440Mod),
	}
}

func (t *TraderConfig) RandomizeParam(min float64, max float64) *TraderConfig {
	idx := randomFloat(0, float64(t.NumParams()))
	switch idx {
	case 0:
		t.BuyPred15Mod = randomFloat(min, max)
		break
	case 1:
		t.BuyPred60Mod = randomFloat(min, max)
		break
	case 2:
		t.BuyPred1440Mod = randomFloat(min, max)
		break
	case 3:
		t.SellPred15Mod = randomFloat(min, max)
		break
	case 4:
		t.SellPred60Mod = randomFloat(min, max)
		break
	case 5:
		t.SellPred1440Mod = randomFloat(min, max)
		break
	case 6:
		t.StopLoss = randomFloat(min, max)
		break
	case 7:
		t.ProfitCap = randomFloat(min, max)
		break
	case 8:
		t.BuyNWQtyMod = randomFloat(min, max)
		break
	case 9:
		t.BuyQty15Mod = randomFloat(min, max)
		break
	case 10:
		t.BuyQty60Mod = randomFloat(min, max)
		break
	case 11:
		t.BuyQty1440Mod = randomFloat(min, max)
		break
	case 12:
		t.SellPosQtyMod = randomFloat(min, max)
		break
	case 13:
		t.SellQty15Mod = randomFloat(min, max)
		break
	case 14:
		t.SellQty60Mod = randomFloat(min, max)
		break
	case 15:
		t.SellQty1440Mod = randomFloat(min, max)
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
