package cursed

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/niljimeno/seamail/config"
)

type Inbox struct {
	Current int
}

func (m model) NewInbox(id int) *Inbox {
	var current int

	for i, mail := range m.Inbox {
		if int(mail.Id) == id {
			current = i
		}
	}

	return &Inbox{
		Current: current,
	}
}

func (s *Inbox) Init() tea.Cmd {
	return getInbox
}

func (s *Inbox) Update(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			s.Current++
			if s.Current >= len(m.Inbox) {
				s.Current = 0
			}

			return m, nil

		case "k", "up":
			s.Current--
			if s.Current < 0 {
				s.Current = len(m.Inbox) - 1
			}

			return m, nil

		case "enter":
			if s.Current >= 0 && s.Current < len(m.Inbox) {
				return m.changeScene(NewSingle(int(m.Inbox[s.Current].Id)))
			}
			return m, tea.Quit
		}

	case updateInbox:
		m.Inbox = msg
		m.Loading = false
		return m, nil
	}
	return m, nil
}

func (s *Inbox) View(m model) string {
	if m.Loading {
		loadingText := lipgloss.NewStyle().
			Render(fmt.Sprintf("Loading messages from %s...", config.Domain))
		return loadingText
	}
	title := lipgloss.NewStyle().
		Foreground(lipgloss.BrightBlue).
		Bold(true).
		Render("Seamail")

	renderedMails := []string{}

	for i, mail := range m.Inbox {
		selected := (i == s.Current)

		var subject string

		if mail.Subject != "" {
			subject = lipgloss.NewStyle().
				Foreground(lipgloss.White).
				Reverse(selected).
				PaddingRight(1).
				Render(mail.Subject)
		}

		address := lipgloss.NewStyle().
			Foreground(lipgloss.Green).
			Reverse(selected).
			PaddingRight(1).
			Render(mail.Address)

		var read string
		if mail.Read {
			read = lipgloss.NewStyle().
				Foreground(lipgloss.Cyan).
				Reverse(selected).
				PaddingRight(1).
				Render("Read")
		} else {
			read = lipgloss.NewStyle().
				Foreground(lipgloss.Red).
				Reverse(selected).
				Bold(true).
				PaddingRight(1).
				Render("Unread")
		}

		mailText := lipgloss.JoinHorizontal(lipgloss.Left, subject, address, read)
		renderedMails = append(
			renderedMails,
			lipgloss.NewStyle().MarginLeft(1).Render(mailText),
		)
	}

	mailBox := lipgloss.JoinVertical(lipgloss.Top, renderedMails...)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		mailBox,
	)
}
