package trader

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader/strategies"
	"strconv"
	"time"
)

var predictions []predictor.Prediction

func SetupEnvironment(startTime time.Time, endTime time.Time, useModel bool, host string, port string) {
	predictions = TrainingData("http://"+host+":"+port+"/aggregator/trader/*", startTime, endTime, useModel)
	log.Println("Locked and Loaded")
}

func TrainingData(serverEndpoint string, startTime time.Time, endTime time.Time, use_model bool) []predictor.Prediction {
	client := http.Client{Timeout: 120 * time.Second}

	var predictions []predictor.Prediction

	req, err := http.NewRequest("GET", serverEndpoint, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("start_time", fmt.Sprint(startTime.Unix()))
	req.Header.Set("end_time", fmt.Sprint(endTime.Unix()))
	req.Header.Set("use_model", fmt.Sprint(use_model))

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		panic("Request Failed : " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&predictions)

	if err != nil {
		panic(err)
	}

	log.Printf("Obtained %d predictions from server...", len(predictions))

	return predictions
}

func RunSingleSim() {
	conf := strategies.BasicConfig{
		BuyPred5Mod:    1.9642530109804408,
		BuyPred10Mod:   0.05497421343571969,
		BuyPred100Mod:  2.4332437134090674,
		SellPred5Mod:   1.3637120887884517,
		SellPred10Mod:  1.238427996702663,
		SellPred100Mod: 2.071777991900559,
		StopLoss:       -0.003641471182833845,
		ProfitCap:      0.02798852232740097,
		BuyQtyMod:      0.04766891255922648,
		SellQtyMod:     0.9980123190092692,
	}

	var _, err = os.Stat("trader.log")
	if !os.IsNotExist(err) {
		os.Remove("trader.log")
	}
	logFile, err := os.OpenFile("trader.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	simulation := NewSimulation(predictions, &conf, 1000, 0.001, 0.1, true)
	simulation.Run()

	fmt.Println(simulation.Trader.Accountant.NetWorth())
}

func RunEvolution() {
	evo := Evolution{
		Predictions:    predictions,
		InitialBalance: 1000,
		Fee:            0.001,
		Uncertainty:    0.5,
		GenerationSize: 200,
		NumGenerations: 15,
		MutationRate:   0.4,
		StartingPoint: []float64{1.9642530109804408, 0.05497421343571969, 2.4332437134090674, 1.3637120887884517, 1.238427996702663,
			2.071777991900559, -0.003641471182833845, 0.02798852232740097, 0.04766891255922648, 0.9980123190092692},
	}

	log.Println("Starting Evo...")

	result := evo.Run()

	log.Println(result.Fitness)
	log.Println(result.Config.ToSlice())
	log.Println("Running single to validate...")

	var _, err = os.Stat("trader.log")
	if !os.IsNotExist(err) {
		os.Remove("trader.log")
	}
	logFile, err := os.OpenFile("trader.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	simulation := NewSimulation(predictions, result.Config, 1000, 0.001, 0, true)
	simulation.Run()

	log.Println(simulation.Trader.Accountant.NetWorth())
}
