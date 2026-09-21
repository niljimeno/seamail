package config

import (
	"bytes"
	"crypto/tls"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/BurntSushi/toml"
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
var Ipv4 bool

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

	Terminal = os.Getenv("TERMINAL")
	Editor = os.Getenv("EDITOR")

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

func LoadCursed() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	seamailConfig := path.Join(configDir, "seamail")
	seamailConfigFile := path.Join(seamailConfig, "config.toml")

	tomlData := struct {
		Domain           string
		Port             int
		SmtpPort         int
		Editor           string
		Terminal         string
		IsTerminalEditor bool
		User             string
		Password         string
		Ipv4             bool
	}{
		Domain:           "example.com",
		Port:             7013,
		SmtpPort:         7012,
		Editor:           os.Getenv("EDITOR"),
		Terminal:         os.Getenv("TERMINAL"),
		IsTerminalEditor: true,
		User:             "nil",
		Password:         "1234",
		Ipv4:             false,
	}

	if _, err := os.Stat(seamailConfigFile); err != nil {
		err = os.MkdirAll(seamailConfig, 0755)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := toml.NewEncoder(&buf).Encode(tomlData); err != nil {
			panic(err)
		}

		os.WriteFile(
			seamailConfigFile,
			buf.Bytes(),
			0755,
		)
	}

	data, err := os.ReadFile(seamailConfigFile)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(data), &tomlData)
	if err != nil {
		return err
	}

	Domain = tomlData.Domain
	ClientPort = tomlData.Port
	AlternativePort = tomlData.SmtpPort
	Editor = tomlData.Editor
	IsTerminalEditor = tomlData.IsTerminalEditor
	Terminal = tomlData.Terminal
	User = tomlData.User
	Password = tomlData.Password
	Ipv4 = tomlData.Ipv4

	if AlternativePort == 0 {
		AlternativePort = 7012
	}

	return nil
}
