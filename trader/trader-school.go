package trader

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
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
	conf := strategies.BasicWithMemoryConfig{
		BuyPred5Mod:    1.064582988619854,
		BuyPred10Mod:   0.7180806459020486,
		BuyPred100Mod:  2.6448109927782526,
		SellPred5Mod:   0.394767696058713,
		SellPred10Mod:  0.5402994113125981,
		SellPred100Mod: 2.344851136724181,
		StopLoss:       -0.003961030174404023,
		ProfitCap:      0.025477934544296355,
		BuyQtyMod:      0.8662148823175331,
		SellQtyMod:     0.9051123877251703,
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
		simulation := NewSimulation(&predictions, strategy, &conf, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001), 0, true)
		simulation.Run()

		fmt.Println(simulation.Trader.Accountant.NetWorth().String() + "$")
	}
	conf2 := strategies.BasicConfig{
		BuyPred5Mod:    1.064582988619854,
		BuyPred10Mod:   0.7180806459020486,
		BuyPred100Mod:  2.6448109927782526,
		SellPred5Mod:   0.394767696058713,
		SellPred10Mod:  0.5402994113125981,
		SellPred100Mod: 2.344851136724181,
		StopLoss:       -0.003961030174404023,
		ProfitCap:      0.025477934544296355,
		BuyQtyMod:      0.8662148823175331,
		SellQtyMod:     0.9051123877251703,
	}
	strategy2 := strategies.NewBasicStrategy(conf2.ToSlice())
	simulation2 := NewSimulation(&predictions, strategy2, &conf2, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001), 0, true)
	simulation2.Run()

	fmt.Println(simulation2.Trader.Accountant.NetWorth().String() + "$  <--- Sem MEM")
}

func RunEvolution() {
	evo := Evolution{
		Predictions:    predictions,
		InitialBalance: decimal.NewFromInt(1000),
		Fee:            decimal.NewFromFloat(0.001),
		Uncertainty:    0,
		GenerationSize: 200,
		NumGenerations: 10,
		MutationRate:   0.4,
		StartingPoint: []float64{1.064582988619854, 0.7180806459020486, 2.6448109927782526, 0.394767696058713,
			0.5402994113125981, 2.344851136724181, -0.003961030174404023, 0.025477934544296355, 0.8662148823175331,
			0.9051123877251703},
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
	simulation := NewSimulation(&predictions, strategy, result.Config, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001), 0, true)
	simulation.Run()

	log.Println(simulation.Trader.Accountant.NetWorth())
}
