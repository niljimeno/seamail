package repository

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-message/mail"
	"github.com/niljimeno/seamail/config"
	"github.com/niljimeno/seamail/models"
)

func Store(data []byte) error {
	id, err := Database.CreateMail([]string{"nil"})
	if err != nil {
		return err
	}

	err = os.WriteFile(
		fmt.Sprintf("%s/mail/%d", config.DataDirectory, id),
		data,
		0755,
	)

	return err
}

func isRead(id []byte) bool {
	_, err := os.Stat(fmt.Sprintf("%s/read/%s", config.DataDirectory, id))
	return !os.IsNotExist(err)
}

func markAsRead(id []byte) {
	os.Create(fmt.Sprintf("%s/read/%s", config.DataDirectory, id))
}

func Get(id []byte) ([]byte, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/mail/%s", config.DataDirectory, id))
	if err != nil {
		return nil, err
	}

	markAsRead(id)
	return content, nil
}

func Remove(id int64) error {
	err := Database.RemoveMail(id, "nil")
	if err != nil {
		log.Println("Could not remove mail")
		return err
	}

	os.Remove(fmt.Sprintf("%s/mail/%d", config.DataDirectory, id))
	return nil
}

func getMailDetails(m models.Mail) (models.MailDetailed, error) {
	result := models.MailDetailed{
		Id:   m.Id,
		Read: m.Read,
	}

	fileDir := fmt.Sprintf("%s/mail/%d", config.DataDirectory, m.Id)
	fileReader, err := os.Open(fileDir)
	if err != nil {
		return result, err
	}

	reader, err := mail.CreateReader(fileReader)
	if err != nil {
		return result, err
	}

	result.Address = reader.Header.Get("From")
	result.Subject = reader.Header.Get("Subject")
	return result, err
}

func ListMail() ([]models.MailDetailed, error) {
	mail, err := Database.ListMail()
	if err != nil {
		log.Println("List mail: Didn't go through the detailed process,", err)
		return nil, err
	}

	result := []models.MailDetailed{}

	for _, m := range mail {
		r, err := getMailDetails(m)
		if err != nil {
			log.Println("Error getting mails:", err)
			continue
		}

		result = append(result, r)
	}

	return result, nil
}
