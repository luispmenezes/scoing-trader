package trader

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"net/http"
	"scoing-trader/trader/model/market"
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
	config := &strategies.BasicWithMemoryConfig{
		BuyPred5Mod:    1.5826542126842869,
		BuyPred10Mod:   2.3353679986593985,
		BuyPred100Mod:  2.3600812220243452,
		SellPred5Mod:   1.1962528097006584,
		SellPred10Mod:  0.5451899361190929,
		SellPred100Mod: 2.9358504210266423,
		StopLoss:       -0.005360676471960717,
		ProfitCap:      0.004446293648514038,
		BuyQtyMod:      0.6436398952398891,
		SellQtyMod:     0.9751410320690478,
	}

	marketEnt := market.NewSimulatedMarket(0, decimal.NewFromFloat(0.001))
	marketEnt.Deposit("USDT", decimal.NewFromInt(1000))

	return &Live{
		HttpClient: http.Client{Timeout: time.Duration(timeout) * time.Second},
		ServerHost: serverHost,
		ServerPort: serverPort,
		Trader: *trader.NewTrader(
			*market.NewAccountant(marketEnt, decimal.NewFromInt(1000), decimal.NewFromFloat(0.001)),
			predictor.NewSimulatedPredictor(0),
			strategies.NewBasicWithMemoryStrategy(config.ToSlice(), 10), true, false),
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
				err := l.Trader.Accountant.UpdateAssetValue(coin, decimal.NewFromFloat(prediction.CloseValue), prediction.Timestamp)
				if err != nil {
					panic(err)
				}
				l.Trader.Predictor.SetNextPrediction(prediction)
				l.Trader.ProcessData(coin)

				if len(l.Trader.Records) != numDecisions {
					for i := int(math.Max(0, float64(numDecisions))); i < len(l.Trader.Records); i++ {
						log.Println(l.Trader.Records[i].ToString())
						numDecisions++
					}
				}
				lastTimestamps[coin] = prediction.Timestamp
			}
		}
		log.Println(l.Trader.Accountant.ToString())
		time.Sleep(60 * time.Second)
	}
}
