package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/niljimeno/seamail/config"
)

func WriteMail() (string, error) {
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
	cmd = exec.Command(config.Editor, tmpFileDir)
	//if config.IsTerminalEditor {
	//	cmd = exec.Command(config.Terminal, config.Editor, tmpFileDir)
	//} else {
	//	cmd = exec.Command(config.Editor, tmpFileDir)
	//}

	if config.IsTerminalEditor {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("using %s which is terminal %s %v : %v - Command: %v", config.Editor, config.Terminal, config.IsTerminalEditor, err, cmd.Path))
	}

	return id, nil
}

func SendMail(id string) error {
	auth := sasl.NewPlainClient("", config.User, config.Password)

	tmpFileDir := fmt.Sprintf("/tmp/mail%s", id)
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

	host := fmt.Sprintf("mail.%s:%d", config.Domain, config.AlternativePort)

	if !config.Ipv4 {
		return smtp.SendMail(host, auth, from, tos, bytes.NewReader(newMsg))
	}

	conn, err := net.Dial("tcp4", host)
	if err != nil {
		conn, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", config.Domain, config.AlternativePort))
		if err != nil {
			return err
		}
	}
	defer conn.Close()

	c, err := smtp.NewClientStartTLS(conn, &tls.Config{ServerName: "mail." + config.Domain})
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Auth(auth); err != nil {
		return err
	}
	if err = c.Mail(from, nil); err != nil {
		return err
	}
	for _, rcpt := range tos {
		if err = c.Rcpt(rcpt, nil); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(newMsg); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return c.Quit()
}
