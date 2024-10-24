package backend

import (
	"log"

	"autocrud/src/codegen"
	"autocrud/src/config"
	"autocrud/src/database"
)

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

func (b BackendGeneratorImpl) Generate() {
	codegen.SetTemplateDir()

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
	if err != nil {
		log.Printf("Error initiating backend module: %v", err)
		return
	}

	config.MakeRelativeDir(b.Directories.Backend, "src")

	filePath := b.Directories.Backend + "/src/main.go"

	err = GenerateMain(filePath, projectName, b.Config)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}

	config.MakeRelativeDir(b.Directories.Backend+"/src", "models")

	modelsDir := b.Directories.Backend + "/src/models/"

	for _, table := range b.Config.Schema.Tables {
		destination := modelsDir + table.Name + ".go"
		err = GenerateModel(destination, table)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
	}

	config.MakeRelativeDir(b.Directories.Backend+"/src", "dao")

	daoDir := b.Directories.Backend + "/src/dao/"
	for _, table := range b.Config.Schema.Tables {
		destination := daoDir + table.Name + "DAO.go"
		daoData := DAOData{
			ProjectName:  projectName,
			Table:        table,
			DatabasePath: b.Directories.DatabasePath,
		}
		err = GenerateDAO(destination, daoData)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
	}

	config.MakeRelativeDir(b.Directories.Backend+"/src", "controller")

	controllerDir := b.Directories.Backend + "/src/controller/"

	destination := controllerDir + "controller.go"
	err = GenerateControllerRouter(destination, projectName)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}

	for _, table := range b.Config.Schema.Tables {
		destination := controllerDir + table.Name + "Controller.go"
		err = GenerateController(destination, projectName, table)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
	}
}
