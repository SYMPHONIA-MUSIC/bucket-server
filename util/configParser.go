package util

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerPort      string `json:"server_port"`
	SecurityEnabled bool   `json:"security_enabled"`
	APIKey          string `json:"api_key"`
	UseHTTPS        bool   `json:"use_https"`
	HTTPSCertPath   string `json:"https_cert_path"`
	HTTPSKeyPath    string `json:"https_key_path"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
