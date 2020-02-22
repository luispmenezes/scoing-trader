package trader

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"scoing-trader/trader/model/market"
	"scoing-trader/trader/model/market/model"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader"
	"scoing-trader/trader/model/trader/strategies"
	"time"
)

type Live struct {
	HttpClient http.Client
	ServerHost string
	ServerPort string
	Trader     trader.Trader
}

const path string = "/predictor/latest/"

var coins = []string{"BTCUSDT", "ETHUSDT", "BNBUSDT", "LTCUSDT", "XRPUSDT"}

func NewLive(serverHost string, serverPort string, timeout int) *Live {
	config := &strategies.BasicConfig{
		BuyPred5Mod:    2.0251376282249756,
		BuyPred10Mod:   0.5873073841927829,
		BuyPred100Mod:  0.5043322053305296,
		SellPred5Mod:   1.2969923715289462,
		SellPred10Mod:  1.16495904843373,
		SellPred100Mod: 0.8470197466391618,
		StopLoss:       -0.2135418301481863,
		ProfitCap:      0.11912187551817349,
		BuyQtyMod:      0.6591540134743608,
		SellQtyMod:     0.6547519264694094,
	}

	marketEnt := market.NewSimulatedMarket(0, 0.001)
	marketEnt.Deposit("USDT", 1000)

	return &Live{
		HttpClient: http.Client{Timeout: time.Duration(timeout) * time.Second},
		ServerHost: serverHost,
		ServerPort: serverPort,
		Trader: *trader.NewTrader(config,
			*market.NewAccountant(marketEnt, model.FloatToInt(1000), 0.001),
			predictor.NewSimulatedPredictor(0),
			strategies.NewBasicStrategy(config.ToSlice()), true),
	}
}

func (l *Live) Run() {
	numDecisions := 0

	lastTimestamps := make(map[string]time.Time)

	log.Println("Starting Live Mode...")

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
					log.Println("Failed getting latest prediction from server.\nSleeping for 30 s...")
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

			lastCoinTimestamp, exists := lastTimestamps[coin]

			if !exists || !prediction.Timestamp.Equal(lastCoinTimestamp) {
				err := l.Trader.Accountant.UpdateAssetValue(coin, model.FloatToInt(prediction.CloseValue), prediction.Timestamp)
				if err != nil {
					panic(err)
				}
				l.Trader.Predictor.SetNextPrediction(prediction)
				l.Trader.ProcessData(coin)

				if len(l.Trader.Records) != numDecisions {
					for i := int(math.Max(0, float64(numDecisions-1))); i < len(l.Trader.Records); i++ {
						log.Println(l.Trader.Records[i].ToString())
					}
					numDecisions = len(l.Trader.Records)
				}
				lastTimestamps[coin] = prediction.Timestamp
			}
		}
		log.Println(l.Trader.Accountant.ToString())
		time.Sleep(60 * time.Second)
	}
}
