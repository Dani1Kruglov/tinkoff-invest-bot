package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type databaseConfig struct {
	DSN string `yaml:"DataBaseDSN"`
}

func LoadConfigDBFile(filename string) (string, error) {
	var dbConfig databaseConfig
	input, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(input, &dbConfig)
	if err != nil {
		log.Println(err)
	}
	return dbConfig.DSN, nil
}
