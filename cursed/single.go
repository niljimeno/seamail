package cursed

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/niljimeno/seamail/models"
)

type Single struct {
	Loading bool
	Id      int
	Mail    models.MailFull
}

type updateSingle models.MailFull

func NewSingle(id int) *Single {
	return &Single{
		Loading: true,
		Id:      id,
	}
}

func (s *Single) Init() tea.Cmd {
	return s.getMessage
}

func (s *Single) getMessage() tea.Msg {
	return getMessage(s.Id)
}

func (s *Single) Update(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateSingle:
		s.Loading = false
		s.Mail = models.MailFull(msg)
		return m, nil
	}
	return m, nil
}

func (s *Single) View(m model) string {
	if s.Loading {
		loadingText := lipgloss.NewStyle().Render("Loading...")
		return loadingText
	}
	return lipgloss.NewStyle().Render("hello")
}
