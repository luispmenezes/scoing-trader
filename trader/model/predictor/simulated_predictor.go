package predictor

import (
	"math/rand"
	"super-trader/trader/model"
)

type SimulatedPredictor struct {
	Predictions map[string][]Prediction
	Index       int
	Uncertainty float64
}

func NewSimulatedPredictor(predictions map[string][]Prediction, uncertainty float64) *SimulatedPredictor {
	return &SimulatedPredictor{
		Predictions: predictions,
		Index:       0,
		Uncertainty: uncertainty,
	}
}

func (p *SimulatedPredictor) Predict(coin string, data model.ExchangeData) float64 {
	prediction := p.Predictions[coin][p.Index].PredictedValue * p.calcError()
	p.Index += 1
	return prediction
}

func (p *SimulatedPredictor) calcError() float64 {
	return 1 - (-p.Uncertainty + rand.Float64()*(2*p.Uncertainty))
}
