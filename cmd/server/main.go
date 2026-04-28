package main

import (
	"log"

	"github.com/niljimeno/seamail/broadcast"
	"github.com/niljimeno/seamail/config"
	"github.com/niljimeno/seamail/repository"
)

func main() {
	log.Println("Loading certs..")
	err := config.LoadServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loading repository..")
	err = repository.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Preparing broadcast..")
	broadcast.Listen()
}
