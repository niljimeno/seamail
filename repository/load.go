package repository

import (
	"fmt"
	"os"

	"github.com/niljimeno/seamail/config"
)

const perms = 0755

func Load() error {
	err := os.MkdirAll(config.DataDirectory, perms)
	if err != nil {
		return err
	}

	os.Mkdir(fmt.Sprintf("%s/mail", config.DataDirectory), perms)
	os.Mkdir(fmt.Sprintf("%s/read", config.DataDirectory), perms)
	os.Mkdir(fmt.Sprintf("%s/banned_addresses", config.DataDirectory), perms)
	os.Mkdir(fmt.Sprintf("%s/banned_domains", config.DataDirectory), perms)

	return nil
}
