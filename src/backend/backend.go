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
	config.MakeRelativeDir(b.Directories.Backend, "src")

	filePath := b.Directories.Backend + "/src/main.go"

	err := codegen.GenerateMain(
		filePath,
		"main",
		"this code was automatically generated",
	)
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

	return nil
}
