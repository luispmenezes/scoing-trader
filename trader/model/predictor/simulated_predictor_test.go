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
		Pred15:     0,
		Pred60:     0,
		Pred1440:   0,
	}, {
		Timestamp:  time.Now().UTC(),
		Coin:       "ETHUSDT",
		CloseValue: 0,
		Pred15:     1,
		Pred60:     1,
		Pred1440:   1,
	}, {
		Timestamp:  time.Now().UTC(),
		Coin:       "BTCUSDT",
		CloseValue: 0,
		Pred15:     2,
		Pred60:     2,
		Pred1440:   2,
	}}

	predictions := make(map[string][]Prediction)
	predictions["BTCUSDT"] = coin_predictions

	predictor := NewSimulatedPredictor(0)
	predictor.SetNextPrediction(coin_predictions[0])
	pred := predictor.Predict("BTCUSDT")
	if pred.Pred15 != 0 || pred.Pred60 != 0 || pred.Pred1440 != 0 {
		t.Error("Expected prediction=0, got ", pred)
	}
	predictor.SetNextPrediction(coin_predictions[1])
	pred = predictor.Predict("ETHUSDT")
	if pred.Pred15 != 1 || pred.Pred60 != 1 || pred.Pred1440 != 1 {
		t.Error("Expected prediction=1, got ", pred)
	}
	predictor.SetNextPrediction(coin_predictions[2])
	pred = predictor.Predict("BTCUSDT")
	if pred.Pred15 != 2 || pred.Pred60 != 2 || pred.Pred1440 != 2 {
		t.Error("Expected prediction=2, got ", pred)
	}
}

func TestUncertainty(t *testing.T) {
	coin_predictions := []Prediction{{
		Timestamp:  time.Now().UTC(),
		Coin:       "BTCUSDT",
		CloseValue: 0,
		Pred15:     1,
		Pred60:     1,
		Pred1440:   1,
	}}

	predictions := make(map[string][]Prediction)
	predictions["BTCUSDT"] = coin_predictions

	predictor := NewSimulatedPredictor(1)
	predictor.SetNextPrediction(coin_predictions[0])
	pred := predictor.Predict("BTCUSDT")
	if pred.Pred15 < -1 || pred.Pred15 > 2 || pred.Pred60 < -1 || pred.Pred60 > 2 || pred.Pred1440 < -1 && pred.Pred1440 > 2 {
		t.Error("Expected prediction between -1 and 2, got ", pred)
	}
}
