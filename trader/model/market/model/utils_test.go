package model

import (
	"fmt"
	"testing"
)

func TestIntFloatMul(t *testing.T) {
	result := IntFloatMul(int64(1234322), 0.5)

	if result != 617161 {
		t.Error(fmt.Sprintf("Multiplication Failed Expected: 617161 Got: %d", result))
	}

	result = IntFloatMul(int64(1234322), 2.0)

	if result != 2468644 {
		t.Error(fmt.Sprintf("Multiplication Failed Expected: 2468644 Got: %d", result))
	}
}

func TestIntFloatDiv(t *testing.T) {
	result := IntFloatDiv(int64(1234322), 2.0)

	if result != 617161 {
		t.Error(fmt.Sprintf("Multiplication Failed Expected: 617161 Got: %d", result))
	}

	result = IntFloatDiv(int64(33333), 2)

	if result != 16666 {
		t.Error(fmt.Sprintf("Multiplication Failed Expected: 16666 Got: %d", result))
	}
}

func TestIntToFloat(t *testing.T) {
	result := IntToFloat(100000000)

	if result != 1.0 {
		t.Error(fmt.Sprintf("Conversion to float failed Expected: 1.0 Got: %f", result))
	}

	result = IntToFloat(1)

	if result != 0.00000001 {
		t.Error(fmt.Sprintf("Conversion to float failed Expected: 0.00000001 Got: %f", result))
	}
}

func TestFloatToInt(t *testing.T) {
	result := FloatToInt(1.0)

	if result != 100000000 {
		t.Error(fmt.Sprintf("Conversion to int failed Expected: 100000000 Got: %d", result))
	}

	result = FloatToInt(0.00000001)

	if result != 1 {
		t.Error(fmt.Sprintf("Conversion to int failed Expected: 1 Got: %d", result))
	}

}

func TestIntToString(t *testing.T) {
	stringVal := IntToString(533330000)

	if stringVal != "5.33330000" {
		t.Error(fmt.Sprintf("ToString failed Expected: 5.33330000 Got: %s", stringVal))
	}

	stringVal = IntToString(33)

	if stringVal != "0.00000033" {
		t.Error(fmt.Sprintf("ToString failed Expected: 0.00000033 Got: %s", stringVal))
	}

	stringVal = IntToString(700000000)

	if stringVal != "7.00000000" {
		t.Error(fmt.Sprintf("ToString failed Expected: 7.00000000 Got: %s", stringVal))
	}
}
