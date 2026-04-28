package client

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/niljimeno/seamail/config"
)

func SendMail() error {
	auth := sasl.NewPlainClient("", config.User, config.Password)

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	msg := []byte("To: \n" +
		"From: " + fmt.Sprintf("%s@%s\n", config.User, config.Domain) +
		"Subject: no_subject\n" +
		"Message-ID: <" + id + fmt.Sprintf("@%s>\n", config.Domain) +
		"\n" +
		"This is the email body.\n")

	tmpFileDir := fmt.Sprintf("/tmp/mail%s", id)
	os.WriteFile(tmpFileDir, msg, 0755)

	var cmd *exec.Cmd
	if config.IsTerminalEditor {
		cmd = exec.Command(config.Terminal, config.Editor, tmpFileDir)
	} else {
		cmd = exec.Command(config.Editor, tmpFileDir)
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Print("Do you want to send the message? [yes/no] ")
	var response string
	fmt.Scan(&response)
	if response != "y" && response != "ye" && response != "yes" {
		return nil
	}

	newMsg, err := os.ReadFile(tmpFileDir)
	if err != nil {
		return err
	}

	msgInfo, err := mail.CreateReader(bytes.NewReader(newMsg))
	if err != nil {
		return err
	}

	from := msgInfo.Header.Get("from")
	to := msgInfo.Header.Get("to")

	tos := strings.Split(to, ", ")

	err = smtp.SendMail(
		fmt.Sprintf("mail.%s:%d", config.Domain, config.AlternativePort),
		auth,
		from,
		tos,
		bytes.NewReader(newMsg),
	)

	if err != nil {
		return err
	}

	return nil
}
