package backend

import (
	"fmt"

	"autocrud/src/codegen"
	"autocrud/src/config"
	"autocrud/src/database"
)

type BackendGenerator interface {
	Generate() error
}

type BackendGeneratorImpl struct {
	Config      config.Config
	Directories database.ProjectDirectories
}

func New(
	config config.Config,
	directories database.ProjectDirectories,
) BackendGeneratorImpl {
	return BackendGeneratorImpl{
		Config:      config,
		Directories: directories,
	}
}

func (b BackendGeneratorImpl) Generate() error {
	projectName := "backend"
	err := config.MultiRunCmdInDir(
		b.Directories.Backend,
		config.Command{
			Cmd:  "go",
			Args: []string{"mod", "init", projectName},
		},
		config.Command{
			Cmd:  "go",
			Args: []string{"get", "github.com/mattn/go-sqlite3"},
		},
		config.Command{
			Cmd:  "go",
			Args: []string{"get", "github.com/gin-gonic/gin"},
		},
		config.Command{
			Cmd:  "go",
			Args: []string{"get", "github.com/gin-contrib/cors"},
		},
	)

	config.MakeRelativeDir(b.Directories.Backend, "src")

	filePath := b.Directories.Backend + "/src/main.go"

	err = codegen.GenerateMain(filePath, projectName, b.Config)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	config.MakeRelativeDir(b.Directories.Backend+"/src", "models")

	modelsDir := b.Directories.Backend + "/src/models/"

	for _, table := range b.Config.Schema.Tables {
		destination := modelsDir + table.Name + ".go"
		err = codegen.GenerateModel(destination, table)
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}
	}

	config.MakeRelativeDir(b.Directories.Backend+"/src", "dao")

	daoDir := b.Directories.Backend + "/src/dao/"
	for _, table := range b.Config.Schema.Tables {
		destination := daoDir + table.Name + "DAO.go"
		daoData := codegen.DAOData{
			ProjectName:  projectName,
			Table:        table,
			DatabasePath: b.Directories.DatabasePath,
		}
		err = codegen.GenerateDAO(destination, daoData)
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}
	}

	config.MakeRelativeDir(b.Directories.Backend+"/src", "controller")

	controllerDir := b.Directories.Backend + "/src/controller/"

	destination := controllerDir + "controller.go"
	err = codegen.GenerateControllerRouter(destination, projectName)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	for _, table := range b.Config.Schema.Tables {
		destination := controllerDir + table.Name + "Controller.go"
		err = codegen.GenerateController(destination, projectName, table)
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}
	}
	return nil
}
