// Package config реализует объект конфигурации, а так же ее загрузку из файла в проект
package config

import (
	"encoding/json"
	"os"
)

// Конфигурация серверов и стратегий
type Config struct {
	Port      string   `json:"port"`     // порт балансировщика
	Backends  []string `json:"backends"` // список адресов бэкенд-серверов
	Strategy  string   `json:"strategy"` // стратегия (round-robin, least_connections, random, по дефолту - round-robin)
	RateLimit struct {
		Capacity   int `json:"capacity"`    // максимальное кол-во токенов
		RefillRate int `json:"refill_rate"` // скорость пополнения токенов
	} `json:"rate_limit"`
}

// LoadConfig загружает конфиг из json-файла
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
