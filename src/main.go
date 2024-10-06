package main

import (
	"log"

	"autocrud/src/backend"
	"autocrud/src/config"
	"autocrud/src/database"
)

func main() {
	conf, err := config.Parse("config_test.yaml")
	if err != nil {
		log.Printf("error reading config: %v", err)
		return
	}

	log.Printf("%v", conf)

	directories, err := database.CreateDbIfNecessary(conf)
	if err != nil {
		log.Printf("%v", conf)
		return
	}

	generator := backend.New(conf, directories)

	err = generator.Generate()
	if err != nil {
		log.Printf("%v", err)
		return
	}

	log.Print("success")
}
