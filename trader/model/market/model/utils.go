package model

import (
	"fmt"
	"math"
)

func IntFloatMul(a int64, b float64) int64 {
	return int64(float64(a) * b)
}

func IntFloatDiv(a int64, b float64) int64 {
	return int64(float64(a) / b)
}

func IntToFloat(a int64) float64 {
	return float64(a) / 100000000
}

func FloatToInt(a float64) int64 {
	return int64(a * 100000000)
}

func IntToString(a int64) string {
	return fmt.Sprintf("%.8f", float64(a)/100000000)
}

func Max(a int64, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a int64, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func TruncateFloat(value float64) float64 {
	output := math.Pow(10, float64(8))
	return float64(round(value*output)) / output
}
