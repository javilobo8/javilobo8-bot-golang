package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	MainChannel string          `json:"mainChannel"`
	Channels    []ChannelConfig `json:"channels"`
}

type ChannelConfig struct {
	Channel        string `json:"channel"`
	Enabled        bool   `json:"enabled"`
	BanCopyWindow  int64  `json:"banCopyWindow"`
	BanCopyEnabled bool   `json:"banCopyEnabled"`
	BanCopyTime    int    `json:"banCopyTime"`
}

const CONFIG_PATH = "./config.json"

var Config Configuration

func ReadConfig() {
	jsonFile, err := os.Open(CONFIG_PATH)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteData, _ := ioutil.ReadAll(jsonFile)
	var config Configuration
	json.Unmarshal(byteData, &config)
	Config = config
}

func WriteConfig() {

}
