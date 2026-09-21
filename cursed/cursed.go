package cursed

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/niljimeno/seamail/models"
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
	Init() tea.Cmd
	Update(tea.Msg, model) (model, tea.Cmd)
	View(model) string
}

type model struct {
	Loading bool
	Scene   Scene
	Inbox   []models.MailDetailed
}

type FullMail struct {
}

func (m model) changeScene(s Scene) (model, tea.Cmd) {
	m.Scene = s
	return m, s.Init()
}

func (m model) Init() tea.Cmd {
	return m.Scene.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "i":
			return m.changeScene(NewWrite())
		}

	case error:
		return m.changeScene(NewError(msg))
	}

	return m.Scene.Update(msg, m)
}

func (m model) View() tea.View {
	v := tea.NewView(m.Scene.View(m))
	v.AltScreen = true
	return v
}
