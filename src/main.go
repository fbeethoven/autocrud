package main

import (
	"log"

	"autocrud/src/config"
	"autocrud/src/database"
)

func main() {
	config, err := config.Parse("config_test.yaml")
	if err != nil {
		log.Printf("error reading config: %v", err)
		return
	}

	log.Printf("%v", config)

	err = database.CreateDbIfNecessary(config)
	if err != nil {
		log.Printf("%v", config)
		return
	}

	log.Print("success")
}
