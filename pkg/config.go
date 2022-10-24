package pkg

import (
	"bufio"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type DBConfig struct {
	Host     string `toml:"db_host"`
	Port     string `toml:"db_port"`
	User     string `toml:"db_user"`
	Password string `toml:"db_password"`
	Name     string `toml:"db_name"`
}

type Config struct {
	ApiKey string `toml:"apiKey"`
	Port   string `toml:"port"`
	*DBConfig
	Cities []string
}

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
