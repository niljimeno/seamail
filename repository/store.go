package repository

import (
	"fmt"
	"io"
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

func RemoveMail(id int64) error {
	err := Database.RemoveMail(id, "nil")
	if err != nil {
		log.Println("Could not remove mail")
		return err
	}

	os.Remove(fmt.Sprintf("%s/mail/%d", config.DataDirectory, id))
	return nil
}

func getMailReader(id int64) (*mail.Reader, error) {
	fileDir := fmt.Sprintf("%s/mail/%d", config.DataDirectory, id)
	fileReader, err := os.Open(fileDir)
	if err != nil {
		return nil, err
	}

	reader, err := mail.CreateReader(fileReader)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func getMailDetails(m models.Mail) (models.MailDetailed, error) {
	result := models.MailDetailed{
		Id:   m.Id,
		Read: m.Read,
	}

	reader, err := getMailReader(m.Id)
	if err != nil {
		return result, err
	}

	result.Address = reader.Header.Get("From")
	result.Subject = reader.Header.Get("Subject")
	return result, err
}

func ReadMail(id int64) (models.MailFull, error) {
	result := models.MailFull{
		Id: id,
	}

	var err error
	result.Read, _ = Database.IsRead(id)

	reader, err := getMailReader(id)
	if err != nil {
		return result, err
	}

	result.Address = reader.Header.Get("From")
	result.Subject = reader.Header.Get("Subject")

	content := []models.MailContent{}
	var c models.MailContent

	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		b, err := io.ReadAll(p.Body)
		if err != nil {
			continue
		}

		c = models.MailContent{}
		c.ContentType = p.Header.Get("Content-Type")
		c.Data = string(b)

		content = append(content, c)
	}

	result.Content = content
	Database.MarkAsRead(id)

	return result, nil
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
