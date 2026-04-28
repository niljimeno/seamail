package broadcast

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/niljimeno/seamail/config"
	"github.com/niljimeno/seamail/repository"
)

type Client struct {
	Connection net.Conn
	Pass       []byte
	Action     []byte
	Args       []byte
}

func handleClient(client *Client) {
	defer client.Connection.Close()
	io.LimitReader(client.Connection, 256)

	data, err := bufio.NewReader(client.Connection).ReadSlice(byte('\r'))
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	log.Printf("Data: %s\n", data)

	lines := bufio.NewScanner(bytes.NewReader(data))
	if !lines.Scan() {
		log.Println("Error: not enough data (0)")
		return
	}

	log.Printf("%s\n", lines.Bytes())
	var found bool
	client.Pass, found = bytes.CutPrefix(lines.Bytes(), []byte("Pass: "))
	if !found {
		log.Println("Error: incorrect format")
		return
	}

	if !lines.Scan() {
		log.Println("Error: not enough data (1)")
		return
	}
	client.Action, found = bytes.CutPrefix(lines.Bytes(), []byte("Action: "))
	if !found {
		log.Println("Error: incorrect format")
		return
	}

	lines.Scan()
	client.Args, found = bytes.CutPrefix(lines.Bytes(), []byte("Args: "))
	if !found {
		log.Println("Error: incorrect format")
		return
	}

	log.Printf("Data recieved :: Pass: %s, Action: %s, Args: %s\n", client.Pass, client.Action, client.Args)
	switch string(client.Action) {
	default:
		client.Connection.Write([]byte("Action not understood"))
	case "list":
		log.Println("Sending: ", repository.List())
		client.Connection.Write([]byte(repository.List()))
	case "get":
		contents, err := repository.Get(client.Args)
		if err != nil {
			log.Println("Error: ", err)
		}
		client.Connection.Write([]byte(contents))
	case "remove":
		err := repository.Remove(client.Args)
		if err != nil {
			log.Println("Error: ", err)
		}
		client.Connection.Write([]byte("removal ok"))
	}
}

func listenClient(port int) {
	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), &config.TlsConfig)
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}

		go handleClient(&Client{Connection: conn})
	}
}
