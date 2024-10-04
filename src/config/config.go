package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var validFields = map[string]int{
	"int":       1,
	"string":    1,
	"timestamp": 1,
}

type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Schema  Schema `yaml:"schema"`
}

type FieldSchema struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type TableSchema struct {
	Name   string        `yaml:"name"`
	Fields []FieldSchema `yaml:"fields"`
}

type Schema struct {
	Tables []TableSchema `yaml:"tables"`
}

func Parse(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	if config.Version != "0.1.0" {
		return Config{}, errors.New("unknown version")
	}

	if err := validateSchema(config.Schema); err != nil {
		return Config{}, err
	}

	return config, nil
}

func validateSchema(schema Schema) error {
	for _, table := range schema.Tables {
		for _, field := range table.Fields {
			if _, ok := validFields[field.Type]; !ok {
				return fmt.Errorf("unknown field: %s", field.Type)
			}
		}
	}
	return nil
}
