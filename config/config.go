package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token string `yaml:"token"`
	Cron  string `yaml:"cron"`
}

func defaultConfig() Config {
	var config = Config{
		Token: "",
		Cron:  "0 0 * * *",
	}
	return config
}

func LoadConfig() (Config, error) {
	fileName := "config.yaml"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		config := defaultConfig()
		data, err := yaml.Marshal(config)
		if err != nil {
			return Config{}, err
		}
		err = os.WriteFile(fileName, data, 0644)
		if err != nil {
			return Config{}, err
		}
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	if config.Token == "" {
		panic(errors.New("token is empty"))
	}
	return config, nil
}
