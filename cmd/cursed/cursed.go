package main

import (
	"log"

	"github.com/niljimeno/seamail/config"
	"github.com/niljimeno/seamail/cursed"
)

func main() {
	err := config.LoadCursed()
	if err != nil {
		log.Println("Error reading configuration:", err)
	}
	cursed.Run()
}
