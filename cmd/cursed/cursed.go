package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"encoding/json/v2"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	p := tea.NewProgram(model{Scene: Inbox, Loading: true})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// scene
const (
	Inbox = 1
)

type model struct {
	Loading bool
	Scene   int
	Inbox   []Mail
}

type Mail struct {
	Id      int    `json:"id"`
	Read    bool   `json:"read"`
	Address string `json:"address"`
	Subject string `json:"subject"`
}

func (m model) Init() tea.Cmd {
	return fetchMails
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		}

	case mailsMsg:
		m.Inbox = msg
		m.Loading = false
		return m, nil
	}
	return m, nil
}

func (m model) View() tea.View {
	var res string
	if m.Loading {
		loadingText := lipgloss.NewStyle().Render("Loading...")
		res = loadingText
	} else {
		res = fmt.Sprintf("loaded %v", m.Inbox)
	}

	return tea.NewView(res)
}

type mailsMsg []Mail

func fetchMails() tea.Msg {
	resp, err := http.Get("http://niliara.net:6000/api/")
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
	return mailsMsg(mails)
}
