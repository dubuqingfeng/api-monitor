package utils

import (
	"github.com/jinzhu/configor"
	"time"
)

// dsn
type MySQLDSN struct {
	Name string
	DSN  string
}

// mysql db
type MySQLDB struct {
	Read     MySQLDSN
	Write    MySQLDSN
	Timezone string
}

// sender config
type SenderConfig struct {
	BearyChat struct {
		IsEnabled      bool
		GroupEndpoint  string
		UnSupportTypes map[string]int
	}
	Slack struct {
		IsEnabled      bool
		RobotToken     string
		Channel        string
		UnSupportTypes map[string]int
	}
	Log struct {
		IsEnabled bool
	}
}

// config
var Config = struct {
	Name                         string `default:"api-monitor"`
	MonitorName                  string `default:"api-monitor"`
	IsDebug                      bool
	MonitoringExpression         string
	APIMonitorEnabled            bool `default:"true"`
	PingAPIMonitorEnabled        bool
	GlobalDatabase               MySQLDB
	APIConfigDatabase            MySQLDB
	APIConfigDatabaseTablePrefix string
	Timeout                      struct {
		Timeout               time.Duration `default:"15"`
		TLSHandshakeTimeout   time.Duration `default:"15"`
		ResponseHeaderTimeout time.Duration `default:"15"`
		ExpectContinueTimeout time.Duration `default:"1"`
	}
	SourceFormat struct {
		Ping     string `default:"mysql"`
		API      string `default:"mysql"`
		Endpoint string `default:"mysql"`
	}
	SenderConfig SenderConfig
}{}

// init config, example: config.example
func InitConfig(files string) {
	err := configor.Load(&Config, files)
	if err != nil {
		panic(err)
	}
}

// TODO check cron expression
func GetMonitoringExpression() string {
	if Config.MonitoringExpression == "" {
		return "* * * * *"
	}
	return Config.MonitoringExpression
}

// get all database configs
func GetAllDatabaseConfigs() map[string]string {
	configs := make(map[string]string)
	AddDatabaseConfig(Config.GlobalDatabase, configs)
	AddDatabaseConfig(Config.APIConfigDatabase, configs)
	return configs
}

// add database config
func AddDatabaseConfig(value MySQLDB, configs map[string]string) {
	if value.Read.DSN != "" && value.Read.Name != "" {
		configs[value.Read.Name] = value.Read.DSN
	}
	if value.Write.DSN != "" && value.Write.Name != "" {
		configs[value.Write.Name] = value.Write.DSN
	}
}
