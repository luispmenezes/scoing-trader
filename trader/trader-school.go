package trader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"super-trader/trader/model"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
	"time"
)

var exchangeData map[string][]model.ExchangeData
var predictions map[string][]predictor.Prediction

func SetupEnvironment(dateStart time.Time,coinCSVs map[string]string)	 {

	exchangeData = make(map[string][]model.ExchangeData)
	predictions = make(map[string][]predictor.Prediction)

	for coin, csvPath := range coinCSVs{
		var coinData []model.ExchangeData
		var coinPredictions []predictor.Prediction

		csv_file, err := os.Open(csvPath)
		if err != nil {
			log.Fatal(err)
		}
		csvReader := csv.NewReader(csv_file)

		_, err = csvReader.Read()
		if err !=nil{
			panic(err)
		}

		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			dataEntry := *model.NewExchangeFromSlice(line)
			if dataEntry.OpenTime.Sub(dateStart) > 0{
				coinData = append(coinData, dataEntry)
				predictedValue, err := strconv.ParseFloat(line[len(line)-1], 64)
				if err != nil{
					panic(err)
				}
				coinPredictions = append(coinPredictions, predictor.Prediction{
					Timestamp:      coinData[len(coinData)-1].OpenTime,
					Coin:           coin,
					PredictedValue: predictedValue,
				})
			}
		}

		sort.Slice(coinData, func(i, j int) bool {
			return coinData[i].OpenTime.Sub(coinData[j].OpenTime) < 0
		})

		sort.Slice(coinPredictions, func(i, j int) bool {
			return coinPredictions[i].Timestamp.Sub(coinPredictions[j].Timestamp) < 0
		})

		exchangeData[coin] = coinData
		predictions[coin] = coinPredictions
	}
	log.Println("Locked and Loaded")
}

func RunSingleSim(){

	/*traderConfig := model.TraderConfig{
		BuyThreshold:      0.005,
		IncreaseThreshold: 0.01,
		SellThreshold:     -0.005,
		MinProfit:         0.05,
		MaxLoss:           0.05,
		PositionSizing:    0.05,
	}
*/
	traderConfig := trader.TraderConfig{
		BuyThreshold:      -0.0020870252250869987,
		IncreaseThreshold:  0.20912728469618153,
		SellThreshold:     -0.22445850767653708,
		MinProfit:         -0.2483787238450663,
		MaxLoss:           -0.3337861564204875,
		PositionSizing:    0.36877999111086635,
	}

	simulation := NewSimulation(exchangeData,predictions,traderConfig,1000,0.001,0.05, true)
	simulation.Run()

	log.Println(simulation.Trader.Wallet.NetWorth())
}

func RunEvolution(){
	evo := Evolution{
		ExchangeData:   exchangeData,
		Predictions:    predictions,
		InitialBalance: 1000,
		Fee:            0.001,
		Uncertainty:    0.05,
		GenerationSize: 100,
		NumGenerations: 5,
		MutationRate:   0.2,
	}
	result := evo.Run()

	log.Println(result)
	log.Println("Running single to validate...")

	simulation := NewSimulation(exchangeData,predictions,result.Config,1000,0.001,0.05, false)
	simulation.Run()

	log.Println(simulation.Trader.Wallet.NetWorth())
}