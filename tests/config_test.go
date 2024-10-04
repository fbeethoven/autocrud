package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"autocrud/src/config"
)

func TestConfig(t *testing.T) {
	expected := config.Config{
		Name:    "TestAPI",
		Version: "0.1.0",
		Schema: config.Schema{
			Tables: []config.TableSchema{
				config.TableSchema{
					Name: "user",
					Fields: []config.FieldSchema{
						config.FieldSchema{
							Name: "name",
							Type: "string",
						},
						config.FieldSchema{
							Name: "age",
							Type: "int",
						},
						config.FieldSchema{
							Name: "created_at",
							Type: "timestamp",
						},
					},
				},
				config.TableSchema{
					Name: "note",
					Fields: []config.FieldSchema{
						config.FieldSchema{
							Name: "title",
							Type: "string",
						},
						config.FieldSchema{
							Name: "body",
							Type: "string",
						},
						config.FieldSchema{
							Name: "created_at",
							Type: "timestamp",
						},
						config.FieldSchema{
							Name: "updated_at",
							Type: "timestamp",
						},
					},
				},
			},
		},
	}

	config, err := config.Parse("data/basic_config.yaml")
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, expected, config)
}

func TestConfigValidation(t *testing.T) {
	_, err := config.Parse("data/unknown_field.yaml")
	t.Log(err)
	if err != nil {
		assert.Error(t, err)
	}
}
