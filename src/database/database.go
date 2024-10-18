package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"autocrud/src/config"
)

type ProjectDirectories struct {
	Root         string
	Database     string
	DatabasePath string
	Backend      string
	Frontend     string
}

func createProjectDir(conf config.Config) (ProjectDirectories, error) {
	projectName := strings.ToLower(conf.Name)
	databasePath, err := filepath.Abs(fmt.Sprintf(
		"./%s/database/%s.db",
		projectName,
		strings.ToLower(conf.Name),
	))
	if err != nil {
		return ProjectDirectories{}, err
	}

	directories := ProjectDirectories{
		Root:         fmt.Sprintf("./%s", projectName),
		Database:     fmt.Sprintf("./%s/database", projectName),
		DatabasePath: databasePath,
		Backend:      fmt.Sprintf("./%s/backend", projectName),
		Frontend:     fmt.Sprintf("./%s/frontend", projectName),
	}

	for _, path := range []string{
		directories.Root,
		directories.Database,
		directories.Backend,
		directories.Frontend,
	} {

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

	db, err := sql.Open("sqlite3", directories.DatabasePath)
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
				newFieldQuery = fmt.Sprintf(
					"%s INTEGER PRIMARY KEY AUTOINCREMENT",
					field.Name,
				)
			} else if field.HasDefault {
				newFieldQuery = fmt.Sprintf(
					"%s INTEGER NOT NULL DEFAULT 0",
					field.Name,
				)
			} else {
				newFieldQuery = fmt.Sprintf("%s INTEGER NOT NULL", field.Name)
			}
		} else if field.Type == config.FieldString {
			if field.HasDefault {
				newFieldQuery = fmt.Sprintf(
					"%s TEXT NOT NULL DEFAULT \"\"",
					field.Name,
				)
			} else {
				newFieldQuery = fmt.Sprintf("%s TEXT NOT NULL", field.Name)
			}
		} else if field.Type == config.FieldTimestamp {
			if field.HasDefault {
				newFieldQuery = fmt.Sprintf(
					"%s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP",
					field.Name,
				)
			} else {
				newFieldQuery = fmt.Sprintf("%s DATETIME NOT NULL", field.Name)
			}
		} else {
			log.Printf("error: found table with invalid type")
		}

		fieldQuery = append(fieldQuery, newFieldQuery)
	}

	query += strings.Join(fieldQuery, ",")
	query += " );"

	return query
}

func GetResourceQuery(table config.TableSchema) string {
	return fmt.Sprintf("SELECT * FROM %s;", table.Name)
}

func GetResourceByIdQuery(table config.TableSchema) string {
	resourceIdName := ""
	for _, field := range table.Fields {
		if field.IsPrimaryKey {
			resourceIdName = field.Name
		}
	}

	if resourceIdName == "" {
		return ""
	}

	return fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=?;",
		table.Name,
		resourceIdName,
	)
}
