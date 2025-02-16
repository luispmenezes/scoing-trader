package strategies

import (
	"math/rand"
)

type BasicConfig struct {
	BuyPred5Mod    float64
	BuyPred10Mod   float64
	BuyPred100Mod  float64
	SellPred5Mod   float64
	SellPred10Mod  float64
	SellPred100Mod float64
	StopLoss       float64
	ProfitCap      float64
	BuyQtyMod      float64
	SellQtyMod     float64
}

func (c *BasicConfig) NumParams() int {
	return 10
}

func (c *BasicConfig) ToSlice() []float64 {
	return []float64{c.BuyPred5Mod, c.BuyPred10Mod, c.BuyPred100Mod, c.SellPred5Mod, c.SellPred10Mod, c.SellPred100Mod,
		c.StopLoss, c.ProfitCap, c.BuyQtyMod, c.SellQtyMod}
}

func (c *BasicConfig) FromSlice(slice []float64) {
	c.BuyPred5Mod = slice[0]
	c.BuyPred10Mod = slice[1]
	c.BuyPred100Mod = slice[2]
	c.SellPred5Mod = slice[3]
	c.SellPred10Mod = slice[4]
	c.SellPred100Mod = slice[5]
	c.StopLoss = slice[6]
	c.ProfitCap = slice[7]
	c.BuyQtyMod = slice[8]
	c.SellQtyMod = slice[9]
}

func (c *BasicConfig) ParamRanges() ([]float64, []float64) {
	var min = make([]float64, c.NumParams())
	var max = make([]float64, c.NumParams())
	//BuyPred5Mod
	min[0] = 0
	max[0] = 3
	//BuyPred10Mod
	min[1] = 0
	max[1] = 3
	//BuyPred100Mod
	min[2] = 0
	max[2] = 3
	//SellPred5Mod
	min[3] = 0
	max[3] = 3
	//SellPred10Mod
	min[4] = 0
	max[4] = 3
	//SellPred100Mod
	min[5] = 0
	max[5] = 3
	//StopLoss
	min[6] = -0.3
	max[6] = 0
	//ProfitCap
	min[7] = 0
	max[7] = 0.2
	//BuyQtyMod
	min[8] = 0
	max[8] = 1
	//SellQtyMod
	min[9] = 0
	max[9] = 1

	return min, max
}

func (c *BasicConfig) RandomFromSlices(a []float64, b []float64) {
	var result = make([]float64, c.NumParams())
	for idx := 0; idx < c.NumParams(); idx++ {
		result[idx] = randomFloat(a[idx], b[idx])
	}
	c.FromSlice(result)
}

func (c *BasicConfig) RandomizeParam() {
	idx := rand.Intn(c.NumParams())
	slice := c.ToSlice()
	min, max := c.ParamRanges()

	slice[idx] = randomFloat(min[idx], max[idx])
	c.FromSlice(slice)
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
