package backend

import (
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
	return nil
}
