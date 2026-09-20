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
	case tea.KeyPressMsg:
		switch msg.String() {
		case "backspace":
			return m.changeScene(m.NewInbox(s.Id))
		}
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
	mail := s.Mail

	var subject string
	if mail.Subject != "" {
		subject = lipgloss.NewStyle().Render(mail.Subject)
	}
	address := lipgloss.NewStyle().Render(mail.Address)

	topBar := lipgloss.JoinHorizontal(lipgloss.Left, subject, address)

	var contentText []string
	for _, c := range mail.Content {
		typeText := lipgloss.NewStyle().Foreground(lipgloss.Red).Render(c.ContentType)
		dataText := lipgloss.NewStyle().Render(c.Data)
		contentText = append(contentText, lipgloss.JoinVertical(lipgloss.Top, typeText, dataText))
	}

	textBlock := lipgloss.JoinVertical(lipgloss.Top, contentText...)

	return lipgloss.JoinVertical(lipgloss.Top, topBar, textBlock)
}
