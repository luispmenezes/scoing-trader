package trader

import (
	"encoding/json"
	"log"
	"net/http"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
	"super-trader/trader/model/trader/strategies"
	"super-trader/trader/model/wallet"
	"time"
)

type Live struct {
	HttpClient http.Client
	ServerHost string
	ServerPort string
	Trader     trader.Trader
}

const path string = "/predictor/latest/"

var coins = []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"}

func NewLive(serverHost string, serverPort string, timeout int) *Live {
	config := &strategies.BasicConfig{
		BuyPred15Mod:    1.7964807132040863,
		BuyPred60Mod:    1.4717291842802593,
		BuyPred1440Mod:  0.6516761024844556,
		SellPred15Mod:   1.7102785724218976,
		SellPred60Mod:   2.500594466058227,
		SellPred1440Mod: 0.8877529307313343,
		StopLoss:        -0.003641471182833845,
		ProfitCap:       0.019435102025411398,
		BuyQtyMod:       0.8333981874622558,
		SellQtyMod:      0.9961800350218821,
	}

	return &Live{
		HttpClient: http.Client{Timeout: time.Duration(timeout) * time.Second},
		ServerHost: serverHost,
		ServerPort: serverPort,
		Trader: *trader.NewTrader(config,
			wallet.NewSimulatedWallet(1000, 0.001),
			predictor.NewSimulatedPredictor(0),
			strategies.NewBasicStrategy(config.ToSlice()), true),
	}
}

func (l *Live) Run() {
	numDecisions := 0

	for {
		for _, coin := range coins {
			endpoint := "http://" + l.ServerHost + ":" + l.ServerPort + path + coin

			req, err := http.NewRequest("GET", endpoint, nil)

			if err != nil {
				panic(err)
			}

			var resp *http.Response

			for {
				resp, err = l.HttpClient.Do(req)
				if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
					log.Println("Collector data request failed sleeping for 30 s...")
					if err == nil {
						log.Println(resp.StatusCode)
					}
					time.Sleep(30 * time.Second)
				} else {
					break
				}
			}

			defer resp.Body.Close()

			var prediction predictor.Prediction

			err = json.NewDecoder(resp.Body).Decode(&prediction)

			if err != nil {
				panic(err)
			}

			l.Trader.Wallet.UpdateCoinValue(coin, prediction.CloseValue, prediction.Timestamp)
			l.Trader.Predictor.SetNextPrediction(prediction)
			l.Trader.ProcessData(coin)

			if len(l.Trader.Records) != numDecisions {
				log.Println(l.Trader.Records[len(l.Trader.Records)-1].ToString())
				numDecisions = len(l.Trader.Records)
			}
		}
		log.Println(l.Trader.Wallet.ToString())
		time.Sleep(60 * time.Second)
	}
}
