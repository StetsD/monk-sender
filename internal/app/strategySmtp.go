package app

import (
	"github.com/stetsd/monk-sender/internal/errorsApp"
	"github.com/stetsd/monk-sender/internal/infrastructure/logger"
	"net/smtp"
)

type StrategySmtp struct {
	auth smtp.Auth
	conf map[string]string
}

func (ss *StrategySmtp) Send(msg *onSendMsg) {
	to := []string{msg.Email}
	body := []byte(
		"To: " + msg.Email + "\r\n" +
			"Subject: " + msg.Title + "\r\n" +
			"\r\n" +
			msg.Description + "\r\n")
	err := smtp.SendMail(ss.conf["host"]+":"+ss.conf["port"], ss.auth, ss.conf["username"], to, body)
	if err != nil {
		logger.Log.Error("email is not sent err: " + err.Error())
	}
}

func (ss *StrategySmtp) Init(conf map[string]string) error {
	ss.conf = conf
	username := conf["username"]
	if username == "" {
		return errorsApp.ErrorApp("smtperr: username is required")
	}
	password := conf["password"]
	if password == "" {
		return errorsApp.ErrorApp("smtperr: password is required")
	}
	host := conf["host"]
	if host == "" {
		return errorsApp.ErrorApp("smtperr: host is required")
	}
	port := conf["port"]
	if port == "" {
		return errorsApp.ErrorApp("smtperr: port is required")
	}

	ss.auth = smtp.PlainAuth("", username, password, host)

	return nil
}
