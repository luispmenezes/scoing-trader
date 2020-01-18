package main

import (
	"log"
	"os"
	"super-trader/trader"
	"time"
)

var logFilePath = "trader.log"

//var server = "menz.dynip.sapo.pt"
var server = "localhost"
var port = "8989"
var logToFile = false
var evolution = false
var liveMode = true

var coinCSVs = map[string]string{"BTCUSDT": "/home/menezes/Documents/training-BTCUSDT.csv"}

func main() {

	var startTime = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	var endTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

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
	if liveMode {
		live := trader.NewLive(server, port, 60)
		live.Run()
	} else {
		trader.SetupEnvironment(startTime, endTime, false, server, port)
		if evolution {
			trader.RunEvolution()
		} else {
			trader.RunSingleSim()
		}
	}
}
