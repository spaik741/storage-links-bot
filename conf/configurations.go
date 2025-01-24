package conf

import (
	"encoding/json"
	"log"
	"os"
)

const fileName = "config.json"

type Configuration struct {
	ApiBot      string `json:"apiBot"`
	StoragePath string `json:"storagePath"`
	BatchSize   int    `json:"batchSize"`
	Token       string `json:"token"`
}

func Compile() *Configuration {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panic("Error opening config file", err)
	}
	var conf Configuration
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		log.Panic("Error parsing config file", err)
	}
	return &conf
}
