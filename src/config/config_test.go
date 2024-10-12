package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	expected := Config{
		Name:    "TestAPI",
		Version: "v0.1.0",
		Schema: Schema{
			Tables: []TableSchema{
				{
					Name: "user",
					Fields: []FieldSchema{
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
					Fields: []FieldSchema{
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

	config, err := Parse("testdata/basic_config.yaml")
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, expected, config)
}

func TestConfigValidation(t *testing.T) {
	testCases := []struct {
		Name         string
		ConfigPath   string
		ErrorMessage string
	}{
		{
			Name:         "TestConfigUnknownField",
			ConfigPath:   "testdata/unknown_field.yaml",
			ErrorMessage: UnknownFieldError,
		},
		{
			Name:         "TestConfigNoPrimaryKeyError",
			ConfigPath:   "testdata/no_primary_key.yaml",
			ErrorMessage: NoPrimaryKeyError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := Parse(testCase.ConfigPath)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), testCase.ErrorMessage)
		})
	}
}
