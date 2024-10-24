package main

import (
	"log"

	"github.com/fbeethoven/autocrud/src/backend"
	"github.com/fbeethoven/autocrud/src/cmd"
	"github.com/fbeethoven/autocrud/src/config"
	"github.com/fbeethoven/autocrud/src/database"
	"github.com/fbeethoven/autocrud/src/frontend"
)

func main() {
	args := cmd.GetArgs()

	if args.GenerateConfig {
		config.Generate()
		return
	}

	conf, err := config.Parse(args.ConfigName)
	if err != nil {
		log.Printf("error reading config: %v", err)
		return
	}

	log.Printf("%+v", conf)

	directories, err := database.CreateDbIfNecessary(conf)
	if err != nil {
		log.Printf("%v", conf)
		return
	}

	backend.New(conf, directories).Generate()
	frontend.New(conf, directories).Generate()

	log.Print("success")
}
