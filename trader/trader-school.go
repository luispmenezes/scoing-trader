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
		BuyPred5Mod:    1.2079495905208983,
		BuyPred10Mod:   1.2314340651251743,
		BuyPred100Mod:  2.639287803446922,
		SellPred5Mod:   0.7310100033627728,
		SellPred10Mod:  2.4236266048303667,
		SellPred100Mod: 0.971248749628451,
		StopLoss:       -0.24754584132282575,
		ProfitCap:      0.09154362564165196,
		BuyQtyMod:      0.4571151261645299,
		SellQtyMod:     0.4373028907049203,
		SegTh:          0.02140330866341518,
		HistSegTh:      0.06870754931936505,
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
	simulation := NewSimulation(&predictions, strategy, &conf, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001), 0, true, false)
	simulation.Run()

	fmt.Println(simulation.Trader.Accountant.NetWorth().String() + "$")

	/*for i := 0; i < 5; i++ {
		simulation := NewSimulation(&predictions, strategy, &conf, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001), 0, true)
		simulation.Run()

		fmt.Println(simulation.Trader.Accountant.NetWorth().String() + "$")
	}
	conf2 := strategies.BasicConfig{
		BuyPred5Mod:    1.5940533413689444,
		BuyPred10Mod:   1.6265337296196787,
		BuyPred100Mod:  2.6448109927782526,
		SellPred5Mod:   1.1962528097006584,
		SellPred10Mod:  2.250716035097317,
		SellPred100Mod: 2.9358504210266423,
		StopLoss:       -0.003961030174404023,
		ProfitCap:      0.010777727359352375,
		BuyQtyMod:      0.7042771619721528,
		SellQtyMod:     0.9751410320690478,
	}
	strategy2 := strategies.NewBasicStrategy(conf2.ToSlice())
	simulation2 := NewSimulation(&predictions, strategy2, &conf2, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001), 0, true)
	simulation2.Run()

	fmt.Println(simulation2.Trader.Accountant.NetWorth().String() + "$  <--- Sem MEM")*/
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
		StartingPoint: []float64{0.9668821395093679, 2.7138169720897705, 2.639287803446922, 1.0940385580770726,
			1.6007641962561916, 0.8169274098545057, -0.22309816719590436, 0.21838016983605293, 0.44771076616692107,
			0.4373028907049203, 0.02140330866341518, 0.16750974746124225},
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

	strategy := strategies.NewBasicWithMemoryStrategy(result.Config.ToSlice(), 5)
	simulation := NewSimulation(&predictions, strategy, result.Config, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001), 0, true, false)
	simulation.Run()

	log.Println(simulation.Trader.Accountant.NetWorth())
}
