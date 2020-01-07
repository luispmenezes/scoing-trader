package main

import (
	"log"
	"os"
	"super-trader/trader"
	"time"
)

var logFilePath = "trader.log"
var logToFile = false
var evolution = true

var coinCSVs = map[string]string{"BTCUSDT": "/home/menezes/Documents/training-BTCUSDT.csv"}

func main() {

	var dateStart = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	if logToFile {
		var _, err = os.Stat(logFilePath)
		if !os.IsNotExist(err) {
			os.Remove(logFilePath)
		}
		logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logFile)
	}
	trader.SetupEnvironment(dateStart, coinCSVs, true)
	if evolution {
		trader.RunEvolution()
	} else {
		trader.RunSingleSim()
	}
}
