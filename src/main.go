package main

import (
	"log"

	"autocrud/src/backend"
	"autocrud/src/config"
	"autocrud/src/database"
	"autocrud/src/frontend"
)

func main() {
	conf, err := config.Parse("config_test.yaml")
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
