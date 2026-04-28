package repository

import (
	"fmt"
	"os"

	"github.com/niljimeno/seamail/config"
)

func Load() error {
	if err := os.MkdirAll(config.DataDirectory, 0755); err != nil {
		return err
	}

	os.Mkdir(fmt.Sprintf("%s/mail", config.DataDirectory), 0755)
	os.Mkdir(fmt.Sprintf("%s/banned_addresses", config.DataDirectory), 0755)
	os.Mkdir(fmt.Sprintf("%s/banned_domains", config.DataDirectory), 0755)

	return nil
}
