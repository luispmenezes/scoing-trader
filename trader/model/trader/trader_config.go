package trader

import (
	"math/rand"
)

type TraderConfig struct {
	BuyThreshold      float64
	IncreaseThreshold float64
	SellThreshold     float64
	MinProfit         float64
	MaxLoss           float64
	PositionSizing    float64
	IncreaseSizing    float64
}

func RandomConfig(min,max float64) *TraderConfig{
	return &TraderConfig{
		BuyThreshold:       randomFloat(min,max),
		IncreaseThreshold:  randomFloat(min,max),
		SellThreshold:      randomFloat(min,max),
		MinProfit:          randomFloat(min,max),
		MaxLoss:            randomFloat(min,max),
		PositionSizing:     randomFloat(min,max),
	}
}

func RandomBetweenTwo(a,b TraderConfig) *TraderConfig{
	return &TraderConfig{
		BuyThreshold:       randomFloat(a.BuyThreshold, b.BuyThreshold),
		IncreaseThreshold:  randomFloat(a.IncreaseThreshold, b.IncreaseThreshold),
		SellThreshold:      randomFloat(a.SellThreshold, b.SellThreshold),
		MinProfit:          randomFloat(a.MinProfit, b.MinProfit),
		MaxLoss:            randomFloat(a.MaxLoss, b.MaxLoss),
		PositionSizing:     randomFloat(a.PositionSizing, b.PositionSizing),
		IncreaseSizing:     randomFloat(a.IncreaseThreshold, b.IncreaseSizing),
	}
}

func (t *TraderConfig) RandomizeParam(idx int, min float64, max float64) *TraderConfig{
	if idx >= 0 && idx < 7 {
		switch idx {
		case 0:
			t.BuyThreshold = randomFloat(min,max)
			break
		case 1:
			t.IncreaseThreshold = randomFloat(min, max)
			break
		case 2:
			t.SellThreshold = randomFloat(min, max)
			break
		case 3:
			t.MinProfit = randomFloat(min, max)
			break
		case 4:
			t.MaxLoss = randomFloat(min, max)
			break
		case 5:
			t.PositionSizing = randomFloat(min, max)
			break
		case 6:
			t.IncreaseSizing = randomFloat(min, max)
			break
		}
		return t
	}else {
		return nil
	}
}

func randomFloat(a,b float64) float64{
	var min,max float64

	if a > b {
		max = a
		min = b
	}else{
		max = b
		min = a
	}

	return min + rand.Float64() * (max - min)
}