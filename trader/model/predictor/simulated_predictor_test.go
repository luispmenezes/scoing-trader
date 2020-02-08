package predictor

import (
	"testing"
	"time"
)

func TestPredict(t *testing.T) {

	coin_predictions := []Prediction{{
		Timestamp:  time.Now().UTC(),
		Coin:       "BTCUSDT",
		CloseValue: 0,
		Pred5:      0,
		Pred10:     0,
		Pred100:    0,
	}, {
		Timestamp:  time.Now().UTC(),
		Coin:       "ETHUSDT",
		CloseValue: 0,
		Pred5:      1,
		Pred10:     1,
		Pred100:    1,
	}, {
		Timestamp:  time.Now().UTC(),
		Coin:       "BTCUSDT",
		CloseValue: 0,
		Pred5:      2,
		Pred10:     2,
		Pred100:    2,
	}}

	predictions := make(map[string][]Prediction)
	predictions["BTCUSDT"] = coin_predictions

	predictor := NewSimulatedPredictor(0)
	predictor.SetNextPrediction(coin_predictions[0])
	pred := predictor.Predict("BTCUSDT")
	if pred.Pred5 != 0 || pred.Pred10 != 0 || pred.Pred100 != 0 {
		t.Error("Expected prediction=0, got ", pred)
	}
	predictor.SetNextPrediction(coin_predictions[1])
	pred = predictor.Predict("ETHUSDT")
	if pred.Pred5 != 1 || pred.Pred10 != 1 || pred.Pred100 != 1 {
		t.Error("Expected prediction=1, got ", pred)
	}
	predictor.SetNextPrediction(coin_predictions[2])
	pred = predictor.Predict("BTCUSDT")
	if pred.Pred5 != 2 || pred.Pred10 != 2 || pred.Pred100 != 2 {
		t.Error("Expected prediction=2, got ", pred)
	}
}

func TestUncertainty(t *testing.T) {
	coin_predictions := []Prediction{{
		Timestamp:  time.Now().UTC(),
		Coin:       "BTCUSDT",
		CloseValue: 0,
		Pred5:      1,
		Pred10:     1,
		Pred100:    1,
	}}

	predictions := make(map[string][]Prediction)
	predictions["BTCUSDT"] = coin_predictions

	predictor := NewSimulatedPredictor(1)
	predictor.SetNextPrediction(coin_predictions[0])
	pred := predictor.Predict("BTCUSDT")
	if pred.Pred5 < -1 || pred.Pred5 > 2 || pred.Pred10 < -1 || pred.Pred10 > 2 || pred.Pred100 < -1 && pred.Pred100 > 2 {
		t.Error("Expected prediction between -1 and 2, got ", pred)
	}
}
