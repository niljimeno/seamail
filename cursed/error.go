package cursed

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ErrorMsg struct {
	Err error
}

func NewError(err error) ErrorMsg {
	return ErrorMsg{
		Err: err,
	}
}

func (s ErrorMsg) Init() tea.Cmd {
	return nil
}

func (s ErrorMsg) Update(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter", "escape":
			return m.changeScene(m.NewInbox(0))
		}
	}

	return m, nil
}

func (s ErrorMsg) View(m model) string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Red).
		Render("Error")

	information := lipgloss.NewStyle().
		Render(fmt.Sprintf("%v", s.Err))

	box := lipgloss.NewStyle().
		Padding(2).
		Render(lipgloss.JoinVertical(lipgloss.Center, title, information))

	return box
}
