package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"autocrud/src/config"
)

type ProjectDirectories struct {
	Root     string
	Database string
	Backend  string
	Frontend string
}

func createProjectDir(conf config.Config) (ProjectDirectories, error) {
	projectName := strings.ToLower(conf.Name)
	directories := ProjectDirectories{
		Root:     fmt.Sprintf("./%s", projectName),
		Database: fmt.Sprintf("./%s/database", projectName),
		Backend:  fmt.Sprintf("./%s/backend", projectName),
		Frontend: fmt.Sprintf("./%s/frontend", projectName),
	}

	for _, path := range []string{
		directories.Root,
		directories.Database,
		directories.Backend,
		directories.Frontend,
	} {

		log.Printf("creating directory: %s", path)
		err := config.MakeDir(path)
		if err != nil {
			return directories, err
		}
	}

	return directories, nil
}

func CreateDbIfNecessary(conf config.Config) (ProjectDirectories, error) {
	directories, err := createProjectDir(conf)
	if err != nil {
		return ProjectDirectories{}, err
	}

	databasePath := fmt.Sprintf("%s/%s.db", directories.Database,
		strings.ToLower(conf.Name))

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return ProjectDirectories{}, err
	}
	defer db.Close()

	for _, table := range conf.Schema.Tables {
		err := createTable(db, table)
		if err != nil {
			return ProjectDirectories{}, err
		}
	}

	return directories, nil
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
