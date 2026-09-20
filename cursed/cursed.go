package cursed

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"encoding/json/v2"

	tea "charm.land/bubbletea/v2"
	"github.com/niljimeno/seamail/config"
)

func Run() {
	p := tea.NewProgram(model{
		Scene:   &Inbox{},
		Loading: true,
	})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type Scene interface {
	Update(tea.Msg, model) (model, tea.Cmd)
	View(model) string
}

type model struct {
	Loading bool
	Scene   Scene
	Inbox   []Mail
}

type Mail struct {
	Id      int    `json:"id"`
	Read    bool   `json:"read"`
	Address string `json:"address"`
	Subject string `json:"subject"`
}

func (m model) Init() tea.Cmd {
	return getInbox
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.Scene.Update(msg, m)
}

func (m model) View() tea.View {
	return tea.NewView(m.Scene.View(m))
}

type updateInbox []Mail

func getInbox() tea.Msg {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/", config.Domain, config.ClientPort))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var mails []Mail
	if err := json.Unmarshal(body, &mails); err != nil {
		return err
	}
	return updateInbox(mails)
}
