package trader

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"super-trader/trader/model"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
	"time"
)

var exchangeData map[string][]model.ExchangeData
var predictions []predictor.Prediction

func SetupEnvironment(startTime time.Time, endTime time.Time, useModel bool, host string, port string) {
	if !useModel {
		predictions = TrainingData("http://"+host+":"+port+"/aggregator/trader/BTCUSDT", startTime, endTime)
	} else {
		log.Println("Not Yet Implemented")
	}
	log.Println("Locked and Loaded")
}

func TrainingData(serverEndpoint string, startTime time.Time, endTime time.Time) []predictor.Prediction {
	client := http.Client{Timeout: 60 * time.Second}

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

	return predictions
}

func RunSingleSim() {
	traderConfig := trader.TraderConfig{
		BuyPred15Mod:    -0.6556657405359857,
		BuyPred60Mod:    -0.10513202884384447,
		BuyPred1440Mod:  -0.5058794369370487,
		SellPred15Mod:   -0.020424043433659413,
		SellPred60Mod:   0.18748561837745112,
		SellPred1440Mod: -0.35233506379660956,
		StopLoss:        -0.07231804326332095,
		ProfitCap:       -0.7538757361595431,
		BuyNWQtyMod:     -0.5126836447623637,
		BuyQty15Mod:     -0.3115006433794979,
		BuyQty60Mod:     -0.5659724173245452,
		BuyQty1440Mod:   -0.24578556392780398,
		SellPosQtyMod:   -0.6652484137333834,
		SellQty15Mod:    -0.03482572222405813,
		SellQty60Mod:    -0.034924346173819794,
		SellQty1440Mod:  0.02887440999597651,
	}

	simulation := NewSimulation(predictions, traderConfig, 1000, 0.001, 0.05, true)
	simulation.Run()

	log.Println(simulation.Trader.Wallet.NetWorth())
}

func RunEvolution() {
	evo := Evolution{
		Predictions:    predictions,
		InitialBalance: 1000,
		Fee:            0.001,
		Uncertainty:    0.03,
		GenerationSize: 300,
		NumGenerations: 10,
		MutationRate:   0.4,
	}
	result := evo.Run()

	log.Println(result)
	log.Println("Running single to validate...")

	simulation := NewSimulation(predictions, result.Config, 1000, 0.001, 0, true)
	simulation.Run()

	log.Println(simulation.Trader.Wallet.NetWorth())
}
