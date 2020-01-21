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
		BuyPred15Mod:    0.40880323236069116,
		BuyPred60Mod:    1.2434503828204961,
		BuyPred1440Mod:  0.40822277813909175,
		SellPred15Mod:   1.2896534389129612,
		SellPred60Mod:   1.0965872751631665,
		SellPred1440Mod: 0.7809285759089554,
		StopLoss:        -0.04304173216430472,
		ProfitCap:       0.006650257809231863,
		BuyQtyMod:       0.059374676309544544,
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
