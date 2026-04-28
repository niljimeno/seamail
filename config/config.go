package config

import (
	"crypto/tls"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Domain string
var User string
var Password string

var DataDirectory string

var Terminal string
var Editor string
var IsTerminalEditor bool

var TlsConfig tls.Config

var fullchain string
var privkey string

var AlternativePort int
var ClientPort int

func loadTLS() error {
	cert, err := tls.LoadX509KeyPair(fullchain, privkey)
	if err != nil {
		return err
	}

	TlsConfig = tls.Config{Certificates: []tls.Certificate{cert}}
	return nil
}

func loadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found, continuing anyways")
	}

	if os.Getenv("IsTerminalEditor") == "true" {
		IsTerminalEditor = true
	}

	Domain = os.Getenv("Domain")
	User = os.Getenv("User")
	Password = os.Getenv("Password")

	DataDirectory = os.Getenv("DataDirectory")

	Terminal = os.Getenv("Terminal")
	Editor = os.Getenv("Editor")

	fullchain = os.Getenv("fullchain")
	privkey = os.Getenv("privkey")

	AlternativePort, err = strconv.Atoi(os.Getenv("AlternativePort"))
	if err != nil {
		return err
	}

	ClientPort, err = strconv.Atoi(os.Getenv("ClientPort"))
	if err != nil {
		return err
	}

	return nil
}

func LoadServer() error {
	if err := loadEnv(); err != nil {
		return err
	}

	if err := loadTLS(); err != nil {
		return err
	}

	return nil
}

func LoadClient() error {
	if err := loadEnv(); err != nil {
		return err
	}

	return nil
}
