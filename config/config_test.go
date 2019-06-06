// +build unit

package config

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	config, err := LoadConfig("../config.json.example")
	if err != nil {
		t.Fatalf("err loading config:\n%v", err)
	}

	if config.Host == "" {
		t.Error("No config")
	}
	if config.Key == "" {
		t.Error("No config")
	}
	if config.Secret == "" {
		t.Error("No config")
	}
}

func TestGetMasterConfig(t *testing.T) {
	config, err := LoadMasterConfig("../config.json.example")
	if err != nil {
		t.Fatalf("err loading master config:\n%v", err)
	}

	if config.Master.Host == "" {
		t.Error("No config")
	}
}
