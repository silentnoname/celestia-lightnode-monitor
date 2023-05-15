package main

import (
	"celestia-lightnode-monitor/pkg/alert"
	"celestia-lightnode-monitor/pkg/config"
	"celestia-lightnode-monitor/pkg/log"
	"celestia-lightnode-monitor/pkg/monitor"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.InitLog()
	log.Log.Info("Start celestia-lightnode-monitor")
	log.Log.Info("Loading config")
	cfg, err := config.LoadConfig("config.toml", ".env")
	if err != nil {
		panic(err)
	}
	log.Log.Info("Standard rpc is " + cfg.Node.StandardRPC)
	log.Log.Info("Light node to check：" + strings.Join(cfg.Node.GatewayAPI, ", "))
	log.Log.Info("Will alert when node balance is less than " + strconv.Itoa(cfg.Node.MinimumBalance) + " utia")
	log.Log.Info("Will check node performance every 5 minutes")
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		log.Log.Info("Start to check node performance")
		nodePerformances, err := monitor.CheckNodes(*cfg)
		if err != nil {
			log.Log.Error("Check performance failed", zap.Error(err))
		}
		log.Log.Info("Start to check and send alert")
		alerterr := alert.SendAlertViaDiscord(*cfg, nodePerformances)
		if alerterr != nil {
			log.Log.Error("Failed to send alert", zap.Error(err))
		}
	}

}
