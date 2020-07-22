package go_utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

// Attachment represents an email attachment.
type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

// 邮件
type Mail struct {
	username, password string
	Body               string
	auth               smtp.Auth
	Attachments        map[string]*Attachment
	To                 []string
	Addr               string
}

func NewMail() *Mail {

	mail := &Mail{username: "zhangning@boe.com.cn", password: "Boe888888"}
	mail.Addr = "mail.boe.com.cn:25"
	mail.Attachments = make(map[string]*Attachment)
	mail.Auth()

	return mail
}

func (s *Mail) Auth() {
	s.auth = &MailAuth{s.username, s.password}
}

func (r *Mail) ParseTemplate(data interface{}, templateFile string) error {

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.Body = buf.String()
	return nil
}

// Bytes returns the mail data
func (m *Mail) ToBytes(cc []string, subject string) []byte {

	buf := bytes.NewBuffer(nil)
	buf.WriteString("From: " + m.username + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")

	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	if len(cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(cc, ",") + "\r\n")
	}

	//fix  Encode
	var coder = base64.StdEncoding
	subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(subject)) + "?="

	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := "f46d043c813270fc6b04c2d223da"

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", "text/html"))
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Attachments) > 0 {

		for _, attachment := range m.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					mime := fmt.Sprintf("Content-Type: %s\r\n", mimetype)
					buf.WriteString(mime)
				} else {
					buf.WriteString("Content-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

// Attach attaches a file.
func (m *Mail) attach(file string, inline bool) error {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(file)

	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}

	return nil
}

// Attach attaches a file.
func (m *Mail) Attach(file string) error {
	return m.attach(file, false)
}

// 发送邮件
func (m *Mail) Send(msg []byte) error {

	c, err := smtp.Dial(m.Addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		host, _, _ := net.SplitHostPort(m.Addr)
		config := &tls.Config{InsecureSkipVerify: true, ServerName: host}

		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	if m.auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(m.auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(m.username); err != nil {
		return err
	}
	for _, addr := range m.To {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
