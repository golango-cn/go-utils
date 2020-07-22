package go_utils

import "net/smtp"

type MailAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &MailAuth{username, password}
}

func (a *MailAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
//   return "LOGIN", []byte{}, nil
  return "LOGIN", []byte(a.username), nil
}

func (a *MailAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}


