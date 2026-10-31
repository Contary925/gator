package config

import (
	"os"
	"encoding/json"
	"fmt"
)

type Config struct {
	DbURL string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() Config {
	homepath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting homepath: %v\n", err)
		return Config{}
	}
	data, err := os.ReadFile(homepath+"/.gatorconfig.json")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return Config{}
	}
	var config Config
	err = json.Unmarshal(data, &config) 
	if err != nil {
		fmt.Printf("Error decoding config: %v\n", err)
		return Config{}
	}
	return config
}

func (config *Config) SetUser(username string) {
	config.Current_user_name = username
	data, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Error marshaling config: %v\n", err)
		return
	}
	homepath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting homepath: %v\n", err)
		return
	}
	err = os.WriteFile(homepath+"/.gatorconfig.json", data, 0644)
	return
}