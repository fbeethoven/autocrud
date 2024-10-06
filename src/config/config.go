package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	FieldInt       string = "int"
	FieldString    string = "string"
	FieldTimestamp string = "timestamp"
)

var validFields = map[string]int{
	FieldInt:       1,
	FieldString:    1,
	FieldTimestamp: 1,
}

type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Schema  Schema `yaml:"schema"`
}

type FieldSchema struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	IsPrimaryKey bool   `yaml:"is_primary_key"`
}

type TableSchema struct {
	Name   string        `yaml:"name"`
	Fields []FieldSchema `yaml:"fields"`
}

type Schema struct {
	Tables []TableSchema `yaml:"tables"`
}

const (
	UnknownFieldError string = "unknown field"
	NoPrimaryKeyError string = "no primary key in table"
)

func MakeDir(projectPath string) error {
	_, err := os.Stat(projectPath)
	if !os.IsNotExist(err) {
		return fmt.Errorf("could not create directory %s", projectPath)
	}

	err = os.Mkdir(projectPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func MakeRelativeDir(parentDir, dirPath string) error {
	directoryPath := fmt.Sprintf("./%s/%s", parentDir, dirPath)

	err := MakeDir(directoryPath)
	if err != nil {
		return err
	}

	return nil
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
		foundPrinmaryKey := false
		for _, field := range table.Fields {
			if _, ok := validFields[field.Type]; !ok {
				return fmt.Errorf("%s: %s", UnknownFieldError, field.Type)
			}
			if field.IsPrimaryKey {
				foundPrinmaryKey = true
			}
		}

		if !foundPrinmaryKey {
			return fmt.Errorf("%s: %s", NoPrimaryKeyError, table.Name)
		}
	}
	return nil
}
