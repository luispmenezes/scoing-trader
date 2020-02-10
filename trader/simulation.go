package trader

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
	"super-trader/trader/model/trader/strategies"
	"super-trader/trader/model/wallet"
)

type Simulation struct {
	Predictions []predictor.Prediction
	Trader      trader.Trader
	Logging     bool
}

func NewSimulation(predictions []predictor.Prediction, config trader.StrategyConfig, initialBalance float64, fee float64,
	uncertainty float64, keepRecords bool) *Simulation {
	return &Simulation{
		Predictions: predictions,
		Trader: *trader.NewTrader(config,
			wallet.NewSimulatedWallet(initialBalance, fee),
			predictor.NewSimulatedPredictor(uncertainty),
			strategies.NewBasicStrategy(config.ToSlice()), keepRecords),
		Logging: keepRecords,
	}
}

func (sim *Simulation) Run() {
	numDecisions := 0
	var history_coin = make(map[string]map[string][]string)
	var history_trader = make(map[string][]string)

	for _, pred := range sim.Predictions {
		sim.Trader.Wallet.UpdateCoinValue(pred.Coin, pred.CloseValue, pred.Timestamp)
		sim.Trader.Predictor.SetNextPrediction(pred)
		sim.Trader.ProcessData(pred.Coin)

		if sim.Logging {
			if len(sim.Trader.Records) != numDecisions {
				log.Println(sim.Trader.Records[len(sim.Trader.Records)-1].ToString())
				numDecisions = len(sim.Trader.Records)

				log.Println(sim.Trader.Wallet.ToString())
			}

			if pred.Timestamp.Minute() == 0 {
				history_trader[pred.Timestamp.Format("2006-01-02 15:04:05")] = []string{
					fmt.Sprintf("%f", sim.Trader.Wallet.GetBalance()), fmt.Sprintf("%f", sim.Trader.Wallet.NetWorth())}

				if _, exists := history_coin[pred.Timestamp.Format("2006-01-02 15:04:05")]; !exists {
					history_coin[pred.Timestamp.Format("2006-01-02 15:04:05")] = make(map[string][]string)
				}

				history_coin[pred.Timestamp.Format("2006-01-02 15:04:05")][pred.Coin] = []string{fmt.Sprintf("%f", pred.CloseValue),
					fmt.Sprintf("%d", len(sim.Trader.Wallet.GetPositions(pred.Coin)))}
			}
		}
	}

	if sim.Logging {
		timestamp_keys := make([]string, 0, len(history_trader))
		for k := range history_trader {
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

		var data = [][]string{headers}

		for _, timestamp := range timestamp_keys {
			var entry = []string{timestamp, history_trader[timestamp][0], history_trader[timestamp][1]}

			for _, coin := range coinList {
				if _, contains := history_coin[timestamp][coin]; contains {
					entry = append(entry, history_coin[timestamp][coin][0])
					entry = append(entry, history_coin[timestamp][coin][1])
				} else {
					entry = append(entry, "0")
					entry = append(entry, "0")
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
