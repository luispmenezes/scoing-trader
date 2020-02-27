package strategies

import (
	"math/rand"
)

type BasicWithMemoryConfig struct {
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
	SegTh          float64
	HistSegTh      float64
}

func (c *BasicWithMemoryConfig) NumParams() int {
	return 12
}

func (c *BasicWithMemoryConfig) ToSlice() []float64 {
	return []float64{c.BuyPred5Mod, c.BuyPred10Mod, c.BuyPred100Mod, c.SellPred5Mod, c.SellPred10Mod, c.SellPred100Mod,
		c.StopLoss, c.ProfitCap, c.BuyQtyMod, c.SellQtyMod, c.SegTh, c.HistSegTh}
}

func (c *BasicWithMemoryConfig) FromSlice(slice []float64) {
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
	c.SegTh = slice[10]
	c.HistSegTh = slice[11]
}

func (c *BasicWithMemoryConfig) ParamRanges() ([]float64, []float64) {
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
	min[6] = -0.4
	max[6] = 0
	//ProfitCap
	min[7] = 0
	max[7] = 0.4
	//BuyQtyMod
	min[8] = 0
	max[8] = 1
	//SellQtyMod
	min[9] = 0
	max[9] = 1
	//Segmentation Threshold
	min[10] = 0
	max[10] = 0.2
	//History Segmentation Threshold
	min[11] = 0
	max[11] = 0.2

	return min, max
}

func (c *BasicWithMemoryConfig) RandomFromSlices(a []float64, b []float64) {
	var result = make([]float64, c.NumParams())
	for idx := 0; idx < c.NumParams(); idx++ {
		result[idx] = randomFloat(a[idx], b[idx])
	}
	c.FromSlice(result)
}

func (c *BasicWithMemoryConfig) RandomizeParam() {
	idx := rand.Intn(c.NumParams())
	slice := c.ToSlice()
	min, max := c.ParamRanges()

	slice[idx] = randomFloat(min[idx], max[idx])
	c.FromSlice(slice)
}
