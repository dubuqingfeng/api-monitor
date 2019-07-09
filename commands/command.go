package main

import (
	"github.com/dubuqingfeng/api-monitor/dbs"
	"github.com/dubuqingfeng/api-monitor/fetchers"
	"github.com/dubuqingfeng/api-monitor/utils"
)

func init() {
	utils.InitConfig("./configs/config.yaml")
	dbs.InitMySQLDB()
}

func main() {
	fetch := fetchers.NewAPIFetcher()
	fetch.Handle()
}
