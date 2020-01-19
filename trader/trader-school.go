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

	return predictions
}

func RunSingleSim() {
	/*conf := strategies.BasicConfig{
		BuyPred15Mod:    1.930762732577754,
		BuyPred60Mod:    1.7868518051534483 ,
		BuyPred1440Mod:  0.05520080799089868,
		SellPred15Mod:   0.243155321218021,
		SellPred60Mod:   1.4920916149331402,
		SellPred1440Mod: 1.1212052597086823,
		StopLoss:        -0.19349762963747624,
		ProfitCap:       0.17719916236926414,
		BuyQtyMod:       0.010617832014426627,
		SellQtyMod:      0.9962488983032237,
	}*/

	conf := strategies.BasicConfig{
		BuyPred15Mod:    0.2834961686425913,
		BuyPred60Mod:    1.6494705884986156,
		BuyPred1440Mod:  0.9858951694300432,
		SellPred15Mod:   0.967851251083437,
		SellPred60Mod:   0.5509856817320783,
		SellPred1440Mod: 1.6617781406885666,
		StopLoss:        -0.2984354025918723,
		ProfitCap:       0.004398177584156121,
		BuyQtyMod:       0.024064165286086042,
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

	simulation := NewSimulation(predictions, &conf, 1000, 0.001, 0.05, true)
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
		NumGenerations: 5,
		MutationRate:   0.4,
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
