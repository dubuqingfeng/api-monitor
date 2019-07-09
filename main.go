package main

import (
	"github.com/dubuqingfeng/api-monitor/dbs"
	"github.com/dubuqingfeng/api-monitor/fetchers"
	"github.com/dubuqingfeng/api-monitor/utils"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func init() {
	utils.InitConfig("./configs/config.yaml")
	log.Info(utils.Config)
	dbs.InitMySQLDB()
}

func main() {
	c := cron.New()
	err := c.AddFunc("0 * * * * *", func() {
		fetch := fetchers.NewAPIFetcher()
		fetch.Handle()
	})
	if err != nil {
		log.Error(err)
	}
	c.Start()
	// does not stop any jobs already running
	defer c.Stop()
	// blocking forever
	select {}
}
