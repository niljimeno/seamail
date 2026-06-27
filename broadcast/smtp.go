package broadcast

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"blitiri.com.ar/go/spf"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/niljimeno/seamail/config"
	"github.com/niljimeno/seamail/models"
	"github.com/niljimeno/seamail/repository"
)

// The Backend implements SMTP server methods.
type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	instance := &Session{}
	instance.ip = c.Conn().RemoteAddr().(*net.TCPAddr).IP
	return instance, nil
}

// A Session is returned after successful login.
type Session struct {
	auth bool
	from string
	to   []string
	ip   net.IP
}

// AuthMechanisms returns a slice of available auth mechanisms; only PLAIN is
// supported in this example.
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators.
func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != config.User || password != config.Password {
			return errors.New("Invalid username or password")
		}
		s.auth = true
		return nil
	}), nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func getHost(address string) string {
	_, host, _ := strings.Cut(address, "@")
	return host
}

func (s *Session) validateSdf(host string) bool {
	result, _ := spf.CheckHostWithSender(s.ip, host, s.from)

	if result == spf.Pass {
		log.Println("Spf verified")
		return true
	}

	log.Println("Spf rejected")
	return false
}

func (s *Session) validateDkim(address string, data []byte) bool {
	log.Println("Validating dkim")
	verifications, err := dkim.Verify(bytes.NewReader(data))
	if err != nil {
		log.Println("Error: ", err)
		return false
	}

	for _, v := range verifications {
		if v.Err == nil {
			log.Println("Valid signature for:", v.Domain)
			return true
		}
	}

	log.Println("Invalid dkim signature")
	return false
}

func (s *Session) validateDomain(host string, data []byte) bool {
	return s.validateSdf(host) || s.validateDkim(host, data)
}

func (s *Session) recieveData(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	host := getHost(s.from)

	if !s.validateDomain(host, data) {
		return smtp.ErrAuthFailed
	}

	err = repository.Store(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) sendData(r io.Reader) error {
	if !s.auth {
		return smtp.ErrAuthFailed
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	addresses := make(map[string]bool)
	for _, v := range s.to {
		address := v
		addressSplit := strings.Split(address, "@")
		if len(addressSplit) == 1 {
			return models.AddressFormattedWrong
		}

		domain := addressSplit[1]
		if addresses[domain] {
			continue
		}
		addresses[domain] = true

		mxs, _ := net.LookupMX(domain)
		if len(mxs) == 0 {
			return models.MXNotFound
		}

		mxHostname := mxs[0]

		err = smtp.SendMail(
			mxHostname.Host+":25",
			nil,
			fmt.Sprintf("%s@%s", config.User, config.Domain),
			s.to, strings.NewReader(string(b)),
		)

		if err != nil {
			log.Println("Error: ", err)
			continue
		}
	}

	return nil
}

func (s *Session) isForMe() bool {
	for _, v := range s.to {
		if getHost(v) == config.Domain {
			return true
		}
	}

	return false
}

func (s *Session) Data(r io.Reader) error {
	if s.isForMe() {
		return s.recieveData(r)
	}

	return s.sendData(r)
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
