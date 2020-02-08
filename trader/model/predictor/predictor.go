package predictor

import "time"

type Predictor interface {
	Predict(coin string) Prediction
	SetNextPrediction(prediction Prediction)
}

type Prediction struct {
	Timestamp  time.Time `json:"open_time"`
	Coin       string    `json:"coin"`
	CloseValue float64   `json:"close_value"`
	Pred5      float64   `json:"pred_5"`
	Pred10     float64   `json:"pred_10"`
	Pred100    float64   `json:"pred_100"`
}
