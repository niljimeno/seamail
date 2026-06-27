package main

import (
	"log"

	"github.com/niljimeno/seamail/repository"
)

func db() error {
	err := repository.Database.Connect()
	if err != nil {
		return err
	}

	id, err := repository.Database.CreateMail([]string{"test_user"})
	if err != nil {
		return err
	}

	log.Println("id was", id)

	return nil

}

func main() {
	err := db()
	log.Println("Test results:", err)
}
