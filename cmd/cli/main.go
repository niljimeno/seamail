package main

import (
	"log"
	"os"

	"github.com/niljimeno/seamail/client"
	"github.com/niljimeno/seamail/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Panicln("Not enough arguments!")
	}

	err := config.LoadClient()
	if err != nil {
		log.Panic(err)
	}

	action := os.Args[1]
	err = client.Process(action)
	if err != nil {
		log.Panic(err)
	}
}
