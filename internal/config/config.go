package config

import (
	"os"
	"fmt"
	"encoding/json"
)

const configFileName = "gatorconfig.json"

type Config struct {
	DatabaseURL    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	return fmt.Sprintf("%s/.config/gator/%s", homeDir, configFileName)
}

func Read() Config {
	
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		fmt.Println("Error: ", err)
		return Config{}
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error: ", err)
		return Config{}
	}
	return config
}

func (c *Config) SetUser(usrN string) {
	c.Username = usrN
	jsonData, err := json.MarshalIndent(c, "", " \t")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = os.WriteFile(getConfigPath(), jsonData, 0666)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}