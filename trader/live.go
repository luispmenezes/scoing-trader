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
	config := trader.TraderConfig{
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

	return &Live{
		HttpClient: http.Client{Timeout: time.Duration(timeout) * time.Second},
		ServerHost: serverHost,
		ServerPort: serverPort,
		Trader: *trader.NewTrader(config,
			wallet.NewSimulatedWallet(1000, 0.001),
			predictor.NewSimulatedPredictor(0),
			strategies.NewBasicStrategy(config), true),
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

			l.Trader.Wallet.UpdateCoinValue(coin, prediction.OpenValue, prediction.Timestamp)
			l.Trader.Predictor.SetNextPrediction(prediction)
			l.Trader.ProcessData(coin)

			if len(l.Trader.Decisions) != numDecisions {
				log.Println(l.Trader.Decisions[len(l.Trader.Decisions)-1])
				numDecisions = len(l.Trader.Decisions)
			}

			log.Printf("%s: %s Predictions:(%f,%f,%f) Value %f$ >> NetWorth: %f Balance: %f", prediction.Timestamp,
				prediction.Coin, prediction.Pred15, prediction.Pred60, prediction.Pred1440, prediction.OpenValue, l.Trader.Wallet.NetWorth(), l.Trader.Wallet.GetBalance())
		}
		time.Sleep(60 * time.Second)
	}
}
