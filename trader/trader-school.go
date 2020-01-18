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
		BuyPred15Mod:   0.7220682253127355,
		BuyPred60Mod:   0.3841515091986546,
		BuyPred1440Mod: 0.2816511062409722,
		StopLoss:       -0.20,
		ProfitCap:      0.05,
		BuyNWQtyMod:    0.7602904103196131,
		BuyQty15Mod:    0.5404493979472732,
		BuyQty60Mod:    0.3094901697193291,
		BuyQty1440Mod:  -0.26127667131565757,
		SellPosQtyMod:  -0.5784497038642875,
		SellQty15Mod:   -0.021812453304933435,
		SellQty60Mod:   -0.04740741550138719,
		SellQty1440Mod: -0.025200193131760328,
	}

	/*traderConfig := trader.TraderConfig{
		BuyPred15Mod:   0.0903586916664672 ,
		BuyPred60Mod:   0.4097697753727636,
		BuyPred1440Mod: -0.5544033630234076,
		StopLoss:       -0.17414676852980995,
		ProfitCap:      0.08900994954754723 ,
		BuyNWQtyMod:    -0.5651479355270681,
		BuyQty15Mod:    -0.14977939642036098,
		BuyQty60Mod:    0.25460656533415454,
		BuyQty1440Mod:  0.03790286534019111,
		SellPosQtyMod:  -0.34214099815266874 ,
		SellQty15Mod:   0.4494399266647915,
		SellQty60Mod:   -0.6980746285559483,
		SellQty1440Mod: -0.5791776194000028,
	}*/

	simulation := NewSimulation(predictions, traderConfig, 1000, 0.001, 0.05, true)
	simulation.Run()

	for _, dec := range simulation.Trader.Decisions {
		log.Println(dec)
	}

	log.Println(simulation.Trader.Wallet.NetWorth())
}

func RunEvolution() {
	evo := Evolution{
		Predictions:    predictions,
		InitialBalance: 1000,
		Fee:            0.001,
		Uncertainty:    0.03,
		GenerationSize: 500,
		NumGenerations: 10,
		MutationRate:   0.4,
	}

	log.Println("Starting Evo...")

	result := evo.Run()

	log.Println(result)
	log.Println("Running single to validate...")

	simulation := NewSimulation(predictions, result.Config, 1000, 0.001, 0, true)
	simulation.Run()

	log.Println(simulation.Trader.Wallet.NetWorth())
}
