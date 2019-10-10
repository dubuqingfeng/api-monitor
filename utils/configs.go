package utils

import (
	"github.com/jinzhu/configor"
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
}

// config
var Config = struct {
	Name                         string `default:"app_name"`
	IsDebug                      bool
	APIMonitorEnabled            bool `default:"true"`
	PingAPIMonitorEnabled        bool
	MonitorName                  string `default:"api-monitor"`
	GlobalDatabase               MySQLDB
	APIConfigDatabase            MySQLDB
	APIConfigDatabaseTablePrefix string
	SenderConfig                 SenderConfig
}{}

// init config, example: config.example
func InitConfig(files string) {
	err := configor.Load(&Config, files)
	if err != nil {
		panic(err)
	}
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
