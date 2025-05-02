package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port      string   `json:"port"`
	Backends  []string `json:"backends"`
	Strategy  string   `json:"strategy"`
	RateLimit struct {
		Capacity   int `json:"capacity"`
		RefillRate int `json:"refill_rate"`
	} `json:"rate_limit"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
