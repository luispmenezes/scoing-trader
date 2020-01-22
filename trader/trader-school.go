package trader

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"super-trader/trader/model"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader/strategies"
	"time"
)

var exchangeData map[string][]model.ExchangeData
var predictions []predictor.Prediction

func SetupEnvironment(startTime time.Time, endTime time.Time, useModel bool, host string, port string) {
	if !useModel {
		predictions = TrainingData("http://"+host+":"+port+"/aggregator/trader/*", startTime, endTime)
	} else {
		log.Println("Not Yet Implemented")
	}
	log.Println("Locked and Loaded")
}

func TrainingData(serverEndpoint string, startTime time.Time, endTime time.Time) []predictor.Prediction {
	client := http.Client{Timeout: 120 * time.Second}

	var predictions []predictor.Prediction

	req, err := http.NewRequest("GET", serverEndpoint, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("start_time", fmt.Sprint(startTime.Unix()))
	req.Header.Set("end_time", fmt.Sprint(endTime.Unix()))

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
		BuyPred15Mod:    1.7964807132040863,
		BuyPred60Mod:    1.4717291842802593,
		BuyPred1440Mod:  0.6516761024844556,
		SellPred15Mod:   1.7102785724218976,
		SellPred60Mod:   2.500594466058227,
		SellPred1440Mod: 1.3906215240890492,
		StopLoss:        -0.003641471182833845,
		ProfitCap:       0.019435102025411398,
		BuyQtyMod:       0.7791674502000079,
		SellQtyMod:      0.9961800350218821,
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

	fmt.Println(simulation.Trader.Wallet.NetWorth())

	for _, record := range simulation.Trader.Records {
		log.Println(record.ToString())
	}
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
		StartingPoint: []float64{1.7964807132040863, 1.4717291842802593, 0.6516761024844556, 1.7102785724218976,
			2.500594466058227, 1.3906215240890492, -0.003641471182833845, 0.019435102025411398, 0.7961493374586381,
			0.9961800350218821},
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

	log.Println(simulation.Trader.Wallet.NetWorth())
}
