package pkg

import (
	"bufio"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// DBConfig is config for database
type DBConfig struct {
	Host     string `toml:"db_host"`
	Port     string `toml:"db_port"`
	User     string `toml:"db_user"`
	Password string `toml:"db_password"`
	Name     string `toml:"db_name"`
}

// Config is main config structure which includes database config
type Config struct {
	APIKey string `toml:"apiKey"`
	Port   string `toml:"port"`
	Cities []string
	*DBConfig
}

// NewConfig is config reader and constructor
func NewConfig() *Config {
	config := &Config{}
	_, err := toml.DecodeFile("configs/config.toml", config)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	config.Cities = readCities()
	return config
}

func readCities() []string {
	file, err := os.Open("configs/city_list")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var cities []string
	read := bufio.NewScanner(file)
	for read.Scan() {
		cities = append(cities, read.Text())
	}
	return cities
}
