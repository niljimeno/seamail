package cursed

import (
	"fmt"
	"io"
	"net/http"

	"encoding/json/v2"

	tea "charm.land/bubbletea/v2"
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("recieved response status: %s", resp.Status)
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
	var mails models.MailFull
	if err := json.Unmarshal(body, &mails); err != nil {
		return err
	}
	return updateSingle(mails)
}
