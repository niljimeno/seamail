package cursed

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"encoding/json/v2"

	tea "charm.land/bubbletea/v2"
	"github.com/emersion/go-message/mail"
	"github.com/niljimeno/seamail/config"
	"github.com/niljimeno/seamail/models"
)

type updateInbox []models.MailDetailed
type updateMessage models.MailFull

func baseUrl() string {
	return fmt.Sprintf("http://%s:%d/api/", config.Domain, config.ClientPort)
}

func getInbox() tea.Msg {
	resp, err := http.Get(baseUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var mails []models.MailDetailed
	if err := json.Unmarshal(body, &mails); err != nil {
		return err
	}
	return updateInbox(mails)
}

func getMessage(id int) tea.Msg {
	resp, err := http.Get(fmt.Sprintf("%smail/%d", baseUrl(), id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	/* format message */
	mr, err := mail.CreateReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	instance := models.MailFull{}
	instance.Address = mr.Header.Get("from")
	instance.Subject = mr.Header.Get("subject")

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		b, _ := io.ReadAll(p.Body)
		contentType := p.Header.Get("Content-Type")

		newContent := models.MailContent{
			ContentType: contentType,
			Data:        string(b),
		}

		instance.Content = append(instance.Content, newContent)
	}

	return updateSingle(instance)
}
