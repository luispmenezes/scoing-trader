package trader

import (
	"encoding/csv"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"os"
	"scoing-trader/trader/model/market"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader"
	"sort"
)

type Simulation struct {
	Predictions *[]predictor.Prediction
	Trader      trader.Trader
	Logging     bool
}

func NewSimulation(predictions *[]predictor.Prediction, strategy trader.Strategy, config trader.StrategyConfig, initialBalance decimal.Decimal, fee decimal.Decimal,
	uncertainty float64, keepRecords bool, keepOnlyTransactions bool) *Simulation {
	marketEnt := market.NewSimulatedMarket(0, fee)
	marketEnt.Deposit("USDT", initialBalance)
	return &Simulation{
		Predictions: predictions,
		Trader: *trader.NewTrader(*market.NewAccountant(marketEnt, initialBalance, fee),
			predictor.NewSimulatedPredictor(uncertainty), strategy, keepRecords, keepOnlyTransactions),
		Logging: keepRecords,
	}
}

func (sim *Simulation) Run() {
	numDecisions := 0
	var historyCoin = make(map[string]map[string][]string)
	var historyTrader = make(map[string][]string)

	for _, pred := range *sim.Predictions {
		err := sim.Trader.Accountant.UpdateAssetValue(pred.Coin, decimal.NewFromFloat(pred.CloseValue), pred.Timestamp)
		if err != nil {
			panic(err)
		}
		sim.Trader.Predictor.SetNextPrediction(pred)
		sim.Trader.ProcessData(pred.Coin)

		if sim.Logging {
			if len(sim.Trader.Records) != numDecisions {
				for i := int(math.Max(0, float64(numDecisions-1))); i < len(sim.Trader.Records); i++ {
					log.Println(sim.Trader.Records[i].ToString())
				}

				numDecisions = len(sim.Trader.Records)

				log.Println(sim.Trader.Accountant.ToString())
			}

			if pred.Timestamp.Minute() == 0 {
				balance, _ := sim.Trader.Accountant.GetBalance().Float64()
				nw, _ := sim.Trader.Accountant.NetWorth().Float64()
				historyTrader[pred.Timestamp.Format("2006-01-02 15:04:05")] = []string{fmt.Sprintf("%.4f", balance), fmt.Sprintf("%.4f", nw)}

				if _, exists := historyCoin[pred.Timestamp.Format("2006-01-02 15:04:05")]; !exists {
					historyCoin[pred.Timestamp.Format("2006-01-02 15:04:05")] = make(map[string][]string)
				}

				assetValue, _ := sim.Trader.Accountant.AssetValue(pred.Coin).Float64()

				historyCoin[pred.Timestamp.Format("2006-01-02 15:04:05")][pred.Coin] = []string{fmt.Sprintf("%.4f", pred.CloseValue),
					fmt.Sprintf("%.4f", assetValue)}
			}
		}

		sim.Trader.Accountant.SyncWithMarket()
	}

	if sim.Logging {
		timestamp_keys := make([]string, 0, len(historyTrader))
		for k := range historyTrader {
			timestamp_keys = append(timestamp_keys, k)
		}
		sort.Strings(timestamp_keys)

		var headers = []string{"Timestamp", "Balance", "Networth"}
		//var coinList []string

		/*for coin := range history_coin[timestamp_keys[len(timestamp_keys)/2]]{
			coinList = append(coinList,coin)
			headers = append(headers,coin)
			headers = append(headers,coin+"_positions")
		}*/

		var coinList = []string{"BTCUSDT", "ETHUSDT", "BNBUSDT", "LTCUSDT", "XRPUSDT"}

		for _, coin := range coinList {
			headers = append(headers, coin)
			headers = append(headers, coin+"_positions")
		}

		var data = [][]string{headers}

		for _, timestamp := range timestamp_keys {
			var entry = []string{timestamp, historyTrader[timestamp][0], historyTrader[timestamp][1]}

			for _, coin := range coinList {
				if _, contains := historyCoin[timestamp][coin]; contains {
					entry = append(entry, historyCoin[timestamp][coin][0])
					entry = append(entry, historyCoin[timestamp][coin][1])
				} else {
					entry = append(entry, "")
					entry = append(entry, "")
				}
			}

			data = append(data, entry)
		}

		file, err := os.Create("result.csv")
		if err != nil {
			panic("Cannot create file")
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		for _, value := range data {
			err := writer.Write(value)
			if err != nil {
				panic("Cannot write to file")
			}
		}
	}

}
