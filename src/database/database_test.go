package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"autocrud/src/config"
)

func TestGetCreateTableQuery(t *testing.T) {
	expected := "CREATE TABLE IF NOT EXISTS user ( " +
		"user_id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"name TEXT," +
		"age INTEGER," +
		"created_at DATETIME DEFAULT CURRENT_TIMESTAMP );"

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

	query := getCreateTableQuery(table)

	assert.Equal(t, query, expected)
}
