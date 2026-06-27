package broadcast

import (
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/niljimeno/seamail/api"
	"github.com/niljimeno/seamail/config"
)

func simpleServer(port int) *smtp.Server {
	b := &Backend{}
	s := smtp.NewServer(b)

	s.Addr = fmt.Sprintf("0.0.0.0:%d", port)
	s.Domain = config.Domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = false

	s.TLSConfig = &config.TlsConfig

	return s
}

func Listen() {
	s1 := simpleServer(25)
	s2 := simpleServer(config.AlternativePort)

	log.Printf("Listening at ports %d, %d and %d\n", 25, config.AlternativePort, config.ClientPort)
	go s1.ListenAndServe()
	go s2.ListenAndServe()
	go listenClient(7013)
	api.Run()
}
