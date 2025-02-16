package predictor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const predictionPath string = "/predictor/predict"

type LivePredictor struct {
	HttpClient http.Client
	Endpoint   string
}

func NewLivePredictor(host string, port string, timeout int) *LivePredictor {
	return &LivePredictor{
		HttpClient: http.Client{Timeout: time.Duration(timeout) * time.Second},
		Endpoint:   "http://" + host + ":" + port + predictionPath,
	}
}

func (p *LivePredictor) Predict(coin string) Prediction {
	requestBody, err := json.Marshal("")

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", p.Endpoint, bytes.NewBuffer(requestBody))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.HttpClient.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		panic("Server request failed")
	}

	defer resp.Body.Close()

	var prediction []Prediction

	err = json.NewDecoder(resp.Body).Decode(&prediction)

	if err != nil {
		panic(err)
	}

	return prediction[0]
}

func (p *LivePredictor) SetNextPrediction(prediction Prediction) {
	panic("Do not use in live mode")
}
