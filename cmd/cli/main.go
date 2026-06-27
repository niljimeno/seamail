package main

import (
	"log"
	"os"

	"github.com/niljimeno/seamail/client"
	"github.com/niljimeno/seamail/config"
)

func main() {
	if len(os.Args) < 2 {
		print("Not enough arguments!")
		return
	}

	err := config.LoadClient()
	if err != nil {
		log.Panic(err)
		return
	}

	action := os.Args[1]
	err = client.Process(action)
	if err != nil {
		log.Panic(err)
	}
}
