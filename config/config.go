package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LibraryDirectory string `json:"libraryDirectory"`
}

// LoadConfig reads configuration from a file.
func LoadConfig(filename string) (*Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	config := &Config{}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(config)
	return config, err
}
