package conf

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Config parts
var (
	Github   *TudouGithubConfig
	Database *TudouDatabaseConfig
	Server   *TudouServerConfig
)

type TudouConfig struct {
	Database TudouDatabaseConfig `json:"database"` // database
	Server   TudouServerConfig   `json:"server"`   // server
	Github   TudouGithubConfig   `json:"github"`   // github
}

type TudouServerConfig struct {
	Port int `json:"port"` // 2358
}

type TudouGithubConfig struct {
	Enabled      bool   `json:"enabled"`      // false
	ClientID     string `json:"clientId"`     // ...
	ClientSecret string `json:"clientSecret"` // ...
}

type TudouDatabaseConfig struct {
	Dialect string `json:"dialect"` // sqlite, mysql, postgres, mssql
	DB      string `json:"db"`      // data/db.sqlite3
}

func getConf() (TudouConfig, error) {
	var conf TudouConfig
	cfgPath := "./config.toml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return conf, fmt.Errorf("ERROR: bad config file - %s", err.Error())
	}

	if _, err := toml.Decode(string(data), &conf); err != nil {
		return conf, fmt.Errorf("ERROR: failed to parse config file - %s", err.Error())
	}

	return conf, nil
}

func Load() {
	conf, _ := getConf()
	Github = &conf.Github
	Database = &conf.Database
	Server = &conf.Server
}
