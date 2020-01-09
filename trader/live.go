package trader

import (
	"encoding/json"
	"log"
	"net/http"
	"super-trader/trader/model"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
	"super-trader/trader/model/wallet"
	"time"
)

type Live struct {
	HttpClient http.Client
	ServerHost string
	ServerPort string
	Trader     trader.Trader
}

const dataPath string = "/collector/data/latest/"

var coins = []string{"BTCUSDT"}

func NewLive(serverHost string, serverPort string, timeout int) *Live {
	return &Live{
		HttpClient: http.Client{Timeout: time.Duration(timeout) * time.Second},
		ServerHost: serverHost,
		ServerPort: serverPort,
		Trader: *trader.NewTrader(trader.TraderConfig{
			BuyThreshold:      0.3752903850331908,
			IncreaseThreshold: -0.012730022134450653,
			SellThreshold:     -0.19780686976755377,
			MinProfit:         -0.7671892206551113,
			MaxLoss:           0.11878695356092953,
			PositionSizing:    0.06713966980593267,
			IncreaseSizing:    -0.02473547420121364,
		},
			wallet.NewSimulatedWallet(1000, 0.001),
			predictor.NewLivePredictor(serverHost, serverPort, 60), true),
	}
}

func (l *Live) Run() {
	for {
		for _, coin := range coins {
			endpoint := "http://" + l.ServerHost + ":" + l.ServerPort + dataPath + coin + "/15"

			log.Println(endpoint)

			req, err := http.NewRequest("GET", endpoint, nil)

			if err != nil {
				panic(err)
			}

			resp, err := l.HttpClient.Do(req)

			if err != nil {
				panic(err)
			}

			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				panic("Server request failed")
			}

			defer resp.Body.Close()

			var xData []model.ExchangeData

			err = json.NewDecoder(resp.Body).Decode(&xData)

			if err != nil {
				panic(err)
			}

			log.Println(xData)

			l.Trader.ProcessData(xData[0], coin)
			time.Sleep(60 * time.Second)
		}
	}
}
