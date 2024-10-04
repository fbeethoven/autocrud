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
				{
					Name: "user",
					Fields: []config.FieldSchema{
						{
							Name:         "user_id",
							Type:         "int",
							IsPrimaryKey: true,
						},
						{
							Name: "name",
							Type: "string",
						},
						{
							Name: "age",
							Type: "int",
						},
						{
							Name: "created_at",
							Type: "timestamp",
						},
					},
				},
				{
					Name: "note",
					Fields: []config.FieldSchema{
						{
							Name:         "note_id",
							Type:         "int",
							IsPrimaryKey: true,
						},
						{
							Name: "title",
							Type: "string",
						},
						{
							Name: "body",
							Type: "string",
						},
						{
							Name: "created_at",
							Type: "timestamp",
						},
						{
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
		assert.Contains(t, err.Error(), config.UnknownFieldError)
	}
}

func TestConfigNoPrimaryKey(t *testing.T) {
	_, err := config.Parse("data/no_primary_key.yaml")
	t.Log(err)
	if err != nil {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), config.NoPrimaryKeyError)
	}
}
