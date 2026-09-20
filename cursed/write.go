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

type Write struct {
	MailId string
	Text   string
}

func (s *Write) Init() tea.Cmd {
	mailId, cmd := WriteMail()
	s.MailId = mailId

	return cmd
}

func (s *Write) Update(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "y":
			err := client.SendMail(s.MailId)
			if err != nil {
				panic(err)
			}

			return m.changeScene(m.NewInbox(0))

		case "n":
			return m.changeScene(m.NewInbox(0))
		}
	case finishMail:
	}

	return m, nil
}

func (s *Write) View(m model) string {
	prompt := lipgloss.NewStyle().
		Render("Do you want to send the message? [yes/no]")

	return prompt
}

type finishMail struct{}

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
			return finishMail{}
		},
	)
}
