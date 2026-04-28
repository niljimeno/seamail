package repository

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/niljimeno/seamail/config"
)

func Store(data []byte) error {
	name := time.Now().UnixMilli()
	err := os.WriteFile(
		fmt.Sprintf("%s/mail/%d", config.DataDirectory, name),
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

func Remove(id []byte) {
	os.Remove(fmt.Sprintf("%s/mail/%s", config.DataDirectory, id))
	os.Remove(fmt.Sprintf("%s/read/%s", config.DataDirectory, id))
}

func List() string {
	directories, _ := os.ReadDir(fmt.Sprintf("%s/mail/", config.DataDirectory))
	names := make([]string, len(directories))

	for i := range directories {
		name := directories[i].Name()
		fileDir := fmt.Sprintf("%s/mail/%s", config.DataDirectory, name)
		fileReader, err := os.Open(fileDir)
		if err != nil {
			names[i] = fmt.Sprintf("Error: %v", err)
			continue
		}

		reader, err := mail.CreateReader(fileReader)
		if err != nil {
			names[i] = fmt.Sprintf("Error: %v", err)
			continue
		}

		var stateText string
		if isRead([]byte(name)) {
			stateText = "[r] "
		}

		names[i] = fmt.Sprintf(
			"%s: %s%s %s",
			name,
			stateText,
			reader.Header.Get("From"),
			reader.Header.Get("Subject"),
		)
	}

	return strings.Join(names, "\n")
}
