package configs

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
)

type loginAuth struct {
	username, password string
}

func loginAuthStartTLS(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

//SendMailForWermProvider create and send a new mail
func SendMailForWermProvider(pathFile string) {
	//server := GetMailConfig()
	fmt.Println(">>>>>>>>>")
	fmt.Println(m.server + " " + m.mail + " " + m.password)
	me := email.NewMessage("this is the sugget", "and this is the bbody email")
	me.From = mail.Address{Name: "From", Address: m.mail}
	me.To = []string{"hmuschett@gmail.com"}

	err := me.Attach(pathFile)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(me)
	//auth := smtp.PlainAuth("", "shop@victoryswitzerland.com", "S12345678v", "smtp.office365.com")
	auth := loginAuthStartTLS(m.mail, m.password)
	if err := email.Send(m.server, auth, me); err != nil {
		log.Fatal(err)
	}
}
