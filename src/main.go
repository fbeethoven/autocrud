package main

import (
	"log"

	"autocrud/src/config"
)

func main() {
	config, err := config.Parse("config.yaml")
	if err != nil {
		log.Printf("error reading config: %v", err)
	}

	log.Printf("%v", config)
}
