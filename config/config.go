package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config ...
type Config struct {
	Host    string
	Key     string
	Secret  string
	Timeout int64
	DB      struct {
		Host     string
		Login    string
		Password string
		Name     string
	}
	Neural struct {
		Iterations int
		Predict    float64
	}
	Strategy struct {
		Profit   float64
		StopLose float64
		Quantity float32
	}
}

// MasterConfig ...
type MasterConfig struct {
	IsDev  bool
	Master Config
	Dev    Config
}

// LoadConfig ...
func LoadConfig(path string) (*Config, error) {
	config, err := LoadMasterConfig(path)
	if err != nil {
		log.Printf("err loading master config:\n%v", err)
		return nil, err
	}

	if config.IsDev {
		return &config.Dev, nil
	}

	return &config.Master, nil
}

// LoadMasterConfig ...
func LoadMasterConfig(path string) (*MasterConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("err opening path %s:\n%v", path, err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := MasterConfig{}
	err = decoder.Decode(&config)

	return &config, err
}
