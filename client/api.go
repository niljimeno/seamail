package client

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/emersion/go-message/mail"
	"github.com/niljimeno/seamail/config"
	"github.com/niljimeno/seamail/models"
)

func runApi(action, args string) ([]byte, error) {
	conn, err := tls.Dial("tcp", fmt.Sprintf("mail.%s:%d", config.Domain, config.ClientPort), nil)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Fprintf(conn, "Pass: %s\nAction: %s\nArgs: %s\r", config.Password, action, args)
	if err != nil {
		return nil, err
	}

	response, err := io.ReadAll(bufio.NewReader(conn))
	if err != nil {
		return nil, err
	}

	conn.Close()
	return response, nil
}

func listMail() error {
	response, err := runApi("list", "")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", response)
	return nil
}

func readMailContent() error {
	if len(os.Args) < 3 {
		return models.TooFewArguments
	}

	args := os.Args[2]
	response, err := runApi("get", args)
	if err != nil {
		return err
	}

	/* format message */
	mr, err := mail.CreateReader(bytes.NewReader(response))
	if err != nil {
		return err
	}
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		b, _ := io.ReadAll(p.Body)
		contentType := p.Header.Get("Content-Type")
		if contentType != "" {
			fmt.Printf("==== %s ====\n", contentType)
		}
		fmt.Printf("%s\n", b)
	}

	return nil
}

func removeMail() error {
	if len(os.Args) < 3 {
		return models.TooFewArguments
	}

	args := os.Args[2]
	response, err := runApi("remove", args)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", response)
	return nil
}

func Process(action string) error {
	switch action {
	default:
		return models.UnrecognisedAction
	case "list":
		return listMail()
	case "get":
		return readMailContent()
	case "remove":
		return removeMail()
	case "send":
		return SendMail()
	}
}
