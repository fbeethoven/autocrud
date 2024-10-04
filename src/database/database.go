package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"autocrud/src/config"
)

func createDirectory(projectPath string) error {
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

func createProjectDir(config config.Config) (projectDirectories, error) {
	projectName := strings.ToLower(config.Name)
	directories := projectDirectories{
		root:     fmt.Sprintf("./%s", projectName),
		database: fmt.Sprintf("./%s/database", projectName),
		backend:  fmt.Sprintf("./%s/backend", projectName),
		frontend: fmt.Sprintf("./%s/frontend", projectName),
	}

	for _, path := range []string{
		directories.root,
		directories.database,
		directories.backend,
		directories.frontend,
	} {

		log.Printf("creating directory: %s", path)
		err := createDirectory(path)
		if err != nil {
			return directories, err
		}
	}

	return directories, nil
}

type projectDirectories struct {
	root     string
	database string
	backend  string
	frontend string
}

func CreateDbIfNecessary(config config.Config) error {
	directories, err := createProjectDir(config)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", directories.database+"/example.db")
	if err != nil {
		return err
	}
	defer db.Close()

	for _, table := range config.Schema.Tables {
		err := createTable(db, table)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTable(db *sql.DB, table config.TableSchema) error {
	query := getCreateTableQuery(table)

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func getCreateTableQuery(table config.TableSchema) string {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( ", table.Name)
	fieldQuery := make([]string, 0, len(table.Fields))

	for _, field := range table.Fields {
		var newFieldQuery string
		if field.Type == config.FieldInt {
			if field.IsPrimaryKey {
				newFieldQuery = fmt.Sprintf("%s INTEGER PRIMARY KEY AUTOINCREMENT", field.Name)
			} else {
				newFieldQuery = fmt.Sprintf("%s INTEGER", field.Name)
			}
		} else if field.Type == config.FieldString {
			newFieldQuery = fmt.Sprintf("%s TEXT", field.Name)
		} else if field.Type == config.FieldTimestamp {
			newFieldQuery = fmt.Sprintf(
				"%s DATETIME DEFAULT CURRENT_TIMESTAMP", field.Name)
		} else {
			log.Printf("error: found table with invalid type")
		}

		fieldQuery = append(fieldQuery, newFieldQuery)
	}

	query += strings.Join(fieldQuery, ",")
	query += " );"

	/*
			query := `CREATE TABLE IF NOT EXISTS users (
		        id INTEGER PRIMARY KEY AUTOINCREMENT,
		        name TEXT NOT NULL,
		        age INTEGER
		    );`
	*/

	return query
}
