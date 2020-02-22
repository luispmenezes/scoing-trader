package trader

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"scoing-trader/trader/model/market/model"
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
	conf := strategies.BasicWithMemoryConfig{
		BuyPred5Mod:    2.6809206329479554,
		BuyPred10Mod:   1.7064990556801447,
		BuyPred100Mod:  0.09743963941547748,
		SellPred5Mod:   2.541286321117902,
		SellPred10Mod:  1.8694827401857466,
		SellPred100Mod: 2.8514328946070377,
		StopLoss:       -0.2090242723175535,
		ProfitCap:      0.10027724036138984,
		BuyQtyMod:      0.05966509880121879,
		SellQtyMod:     0.3207395712639675,
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
	strategy := strategies.NewBasicWithMemoryStrategy(conf.ToSlice(), 10)

	for i := 0; i < 5; i++ {
		simulation := NewSimulation(&predictions, strategy, &conf, 1000, 0.001, 0, true)
		simulation.Run()

		fmt.Println(model.IntToString(simulation.Trader.Accountant.NetWorth()) + "$")
	}
	conf2 := strategies.BasicConfig{
		BuyPred5Mod:    2.6809206329479554,
		BuyPred10Mod:   1.7064990556801447,
		BuyPred100Mod:  0.09743963941547748,
		SellPred5Mod:   2.541286321117902,
		SellPred10Mod:  1.8694827401857466,
		SellPred100Mod: 2.8514328946070377,
		StopLoss:       -0.2090242723175535,
		ProfitCap:      0.10027724036138984,
		BuyQtyMod:      0.05966509880121879,
		SellQtyMod:     0.3207395712639675,
	}
	strategy2 := strategies.NewBasicStrategy(conf2.ToSlice())
	simulation2 := NewSimulation(&predictions, strategy2, &conf2, 1000, 0.001, 0, true)
	simulation2.Run()

	fmt.Println(model.IntToString(simulation2.Trader.Accountant.NetWorth()) + "$  <--- Sem MEM")
}

func RunEvolution() {
	evo := Evolution{
		Predictions:    predictions,
		InitialBalance: 1000,
		Fee:            0.001,
		Uncertainty:    0,
		GenerationSize: 400,
		NumGenerations: 5,
		MutationRate:   0.4,
		StartingPoint:  nil,
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

	strategy := strategies.NewBasicWithMemoryStrategy(result.Config.ToSlice(), 10)
	simulation := NewSimulation(&predictions, strategy, result.Config, 1000, 0.001, 0, true)
	simulation.Run()

	log.Println(simulation.Trader.Accountant.NetWorth())
}
