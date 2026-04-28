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

func Get(id []byte) ([]byte, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/mail/%s", config.DataDirectory, id))
	if err != nil {
		return nil, err
	}

	return content, nil
}

func Remove(id []byte) error {
	err := os.Remove(fmt.Sprintf("%s/%s", config.DataDirectory, id))
	return err
}

func List() string {
	directories, _ := os.ReadDir(config.DataDirectory)
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

		names[i] = fmt.Sprintf(
			"%s: %s %s",
			name, reader.Header.Get("From"),
			reader.Header.Get("Subject"),
		)
	}

	return strings.Join(names, "\n")
}
