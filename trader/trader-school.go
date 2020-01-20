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
		BuyPred15Mod:    0.40880323236069116,
		BuyPred60Mod:    1.4454241616461232,
		BuyPred1440Mod:  1.044787863906082,
		SellPred15Mod:   1.1119926745396007,
		SellPred60Mod:   1.0762199758255402,
		SellPred1440Mod: 0.5849382153611731,
		StopLoss:        -0.020285351376654492,
		ProfitCap:       0.006650257809231863,
		BuyQtyMod:       0.027574445407428833,
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
		Uncertainty:    0.05,
		GenerationSize: 200,
		NumGenerations: 10,
		MutationRate:   0.4,
		StartingPoint: []float64{0.40880323236069116, 1.4454241616461232, 1.044787863906082, 1.1119926745396007,
			1.0762199758255402, 0.5849382153611731, -0.020285351376654492, 0.006650257809231863, 0.027574445407428833,
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
