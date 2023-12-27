package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

// Print prints all the configuration settings.
func (c *Config) Print() {
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()
	fmt.Println("--------------------------------------------------------")
	fmt.Println("\t\t\tConfiguration")
	fmt.Println("--------------------------------------------------------")
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		name := typ.Field(i).Name
		fmt.Printf("%s: %v\n", name, field.Interface())
	}
	fmt.Println("--------------------------------------------------------")
}
