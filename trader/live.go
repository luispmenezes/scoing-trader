package trader

import (
	"encoding/json"
	"log"
	"net/http"
	"super-trader/trader/model"
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

const dataPath string = "/collector/data/latest/"

var coins = []string{"BTCUSDT"}

func NewLive(serverHost string, serverPort string, timeout int) *Live {
	config := trader.TraderConfig{
		BuyPred15Mod:    -0.6556657405359857,
		BuyPred60Mod:    -0.10513202884384447,
		BuyPred1440Mod:  -0.5058794369370487,
		SellPred15Mod:   -0.020424043433659413,
		SellPred60Mod:   0.18748561837745112,
		SellPred1440Mod: -0.35233506379660956,
		StopLoss:        -0.07231804326332095,
		ProfitCap:       -0.7538757361595431,
		BuyNWQtyMod:     -0.5126836447623637,
		BuyQty15Mod:     -0.3115006433794979,
		BuyQty60Mod:     -0.5659724173245452,
		BuyQty1440Mod:   -0.24578556392780398,
		SellPosQtyMod:   -0.6652484137333834,
		SellQty15Mod:    -0.03482572222405813,
		SellQty60Mod:    -0.034924346173819794,
		SellQty1440Mod:  0.02887440999597651,
	}

	return &Live{
		HttpClient: http.Client{Timeout: time.Duration(timeout) * time.Second},
		ServerHost: serverHost,
		ServerPort: serverPort,
		Trader: *trader.NewTrader(config,
			wallet.NewSimulatedWallet(1000, 0.001),
			predictor.NewLivePredictor(serverHost, serverPort, 60),
			strategies.NewBasicStrategy(config)),
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

			var resp *http.Response

			for {
				resp, err = l.HttpClient.Do(req)
				if err != nil {
					log.Println("Collector data request failed sleeping for 30 s...")
					time.Sleep(30 * time.Second)
				} else {
					break
				}
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

			l.Trader.ProcessData(coin)
			time.Sleep(60 * time.Second)
		}
	}
}
