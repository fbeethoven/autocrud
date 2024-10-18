package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"autocrud/src/config"
)

func TestGetCreateTableQuery(t *testing.T) {
	expected := "CREATE TABLE IF NOT EXISTS user ( " +
		"user_id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"name TEXT NOT NULL," +
		"age INTEGER NOT NULL," +
		"created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP );"

	table := config.TableSchema{
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
				Name:       "created_at",
				Type:       "timestamp",
				HasDefault: true,
			},
		},
	}

	query := getCreateTableQuery(table)

	assert.Equal(t, expected, query)
}

func TestGetResourceQuery(t *testing.T) {
	table := config.TableSchema{
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
	}

	expected := "SELECT * FROM user;"

	assert.Equal(t, expected, GetResourceQuery(table))
}

func TestGetResourceByIdQuery(t *testing.T) {
	table := config.TableSchema{
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
	}

	expected := "SELECT * FROM user WHERE user_id=?;"

	assert.Equal(t, expected, GetResourceByIdQuery(table))
}

func TestConfigWithDefaults(t *testing.T) {
	expected := "CREATE TABLE IF NOT EXISTS user ( " +
		"user_id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"name TEXT NOT NULL DEFAULT \"\"," +
		"age INTEGER NOT NULL DEFAULT 0," +
		"created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP );"

	table := config.TableSchema{
		Name: "user",
		Fields: []config.FieldSchema{
			{
				Name:         "user_id",
				Type:         "int",
				IsPrimaryKey: true,
				HasDefault:   true,
			},
			{
				Name:       "name",
				Type:       "string",
				HasDefault: true,
			},
			{
				Name:       "age",
				Type:       "int",
				HasDefault: true,
			},
			{
				Name:       "created_at",
				Type:       "timestamp",
				HasDefault: true,
			},
		},
	}

	query := getCreateTableQuery(table)

	assert.Equal(t, expected, query)
}
