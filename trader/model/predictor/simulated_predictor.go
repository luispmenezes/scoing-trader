package predictor

import (
	"math/rand"
)

type SimulatedPredictor struct {
	NextPrediction Prediction
	Uncertainty    float64
}

func NewSimulatedPredictor(uncertainty float64) *SimulatedPredictor {
	return &SimulatedPredictor{
		NextPrediction: Prediction{},
		Uncertainty:    uncertainty,
	}
}

func (p *SimulatedPredictor) Predict(coin string) Prediction {
	if coin != p.NextPrediction.Coin {
		panic("Prediction coin: " + p.NextPrediction.Coin + " doesnt match " + coin)
		return Prediction{}
	} else {
		return p.NextPrediction
	}
}

func (p *SimulatedPredictor) SetNextPrediction(prediction Prediction) {
	p.NextPrediction = prediction
	p.NextPrediction.Pred15 *= p.calcError()
	p.NextPrediction.Pred60 *= p.calcError()
	p.NextPrediction.Pred1440 *= p.calcError()
}

func (p *SimulatedPredictor) calcError() float64 {
	return 1 - (-p.Uncertainty + rand.Float64()*(2*p.Uncertainty))
}
