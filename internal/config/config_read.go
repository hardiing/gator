package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Url      string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("url: %s, username: %s\n", cfg.Url, cfg.Username)

	return cfg, err
}

func (cfg *Config) SetUser(userName string) error {
	cfg.Username = userName
	write(*cfg)
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(home, configFileName)
	return configFilePath, nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(home, configFileName)
	err = os.WriteFile(configFilePath, jsonData, 0644)
	if err != nil {
		return err
	}
	return err
}
