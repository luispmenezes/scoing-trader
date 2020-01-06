package predictor

import (
	"super-trader/trader/model"
	"testing"
	"time"
)

func TestPredict(t *testing.T) {

	coin_predictions := []Prediction{{
		Timestamp:      time.Now().UTC(),
		Coin:           "BTCUSDT",
		PredictedValue: 0,
	}, {
		Timestamp:      time.Now().UTC(),
		Coin:           "BTCUSDT",
		PredictedValue: 1,
	}, {
		Timestamp:      time.Now().UTC(),
		Coin:           "BTCUSDT",
		PredictedValue: 2,
	}}

	predictions := make(map[string][]Prediction)
	predictions["BTCUSDT"] = coin_predictions

	predictor := NewSimulatedPredictor(predictions,0)
	pred := predictor.Predict("BTCUSDT",model.ExchangeData{})
	if pred != 0 {
		t.Error("Expected prediction=0, got ", pred)
	}
	pred = predictor.Predict("BTCUSDT",model.ExchangeData{})
	if pred != 1 {
		t.Error("Expected prediction=1, got ", pred)
	}
	pred = predictor.Predict("BTCUSDT",model.ExchangeData{})
	if pred != 2 {
		t.Error("Expected prediction=2, got ", pred)
	}
}

func TestUncertainty(t *testing.T) {
	coin_predictions := []Prediction{{
		Timestamp:      time.Now().UTC(),
		Coin:           "BTCUSDT",
		PredictedValue: 1,
	}}

	predictions := make(map[string][]Prediction)
	predictions["BTCUSDT"] = coin_predictions

	predictor := NewSimulatedPredictor(predictions,1)
	pred := predictor.Predict("BTCUSDT",model.ExchangeData{})
	if pred > -1 && pred < 2{
		t.Error("Expected prediction between -1 and 2, got ", pred)
	}
}