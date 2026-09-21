package cursed

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	// "github.com/niljimeno/seamail/client"
	"github.com/niljimeno/seamail/client"
	"github.com/niljimeno/seamail/config"
)

const (
	Writing      = 1
	Sending      = 2
	Confirmation = 3
)

type Write struct {
	MailId string
	State  int
}

func NewWrite() *Write {
	return &Write{
		State: Writing,
	}
}

func (s *Write) Init() tea.Cmd {
	mailId, cmd := WriteMail()
	s.MailId = mailId

	return cmd
}

func (s *Write) Update(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if s.State == Confirmation {
			return m.changeScene(m.NewInbox(0))
		}
		switch msg.String() {
		case "y":
			s.State = Sending
			return m, s.sendMail

		case "n":
			return m.changeScene(m.NewInbox(0))

		}
	case mailSent:
		s.State = Confirmation
		return m, nil
	}

	return m, nil
}

type mailSent struct{}

func (s *Write) sendMail() tea.Msg {
	err := client.SendMail(s.MailId)
	if err != nil {
		return fmt.Errorf("Could not send mail - %v", err)
	}

	return mailSent{}
}

func (s *Write) View(m model) string {
	switch s.State {
	case Writing:
		prompt := lipgloss.NewStyle().
			Render("Do you want to send the message? [yes/no]")
		return prompt

	case Sending:
		return lipgloss.NewStyle().
			Render(fmt.Sprintf("Sending mail to %s...", config.Domain))

	case Confirmation:
		prompt := lipgloss.NewStyle().
			Render("Message sent!")
		return prompt
	}

	return ""
}

func WriteMail() (string, tea.Cmd) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	msg := []byte("To: \n" +
		"From: " + fmt.Sprintf("%s@%s\n", config.User, config.Domain) +
		"Subject: no_subject\n" +
		"Message-ID: <" + id + fmt.Sprintf("@%s>\n", config.Domain) +
		"\n" +
		"This is the email body.\n")

	tmpFileDir := fmt.Sprintf("/tmp/mail%s", id)
	os.WriteFile(tmpFileDir, msg, 0755)

	return id, tea.ExecProcess(
		exec.Command(config.Editor, tmpFileDir),
		func(err error) tea.Msg {
			return nil
		},
	)
}
