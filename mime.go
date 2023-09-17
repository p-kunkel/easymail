package easymail

import (
	"bytes"
	"mime"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
)

type Mail struct {
	Headers *Header
	Parts   []PartCreator
}

type MIME []byte

type PartCreator interface {
	GetHeaders() textproto.MIMEHeader
	CreateAndWritePartTo(*multipart.Writer) error
}

func New() *Mail {
	return &Mail{
		Headers: &Header{
			encoder: mime.BEncoding,
			charset: "UTF-8",
			custom:  textproto.MIMEHeader{},
		},
	}
}

//Returns a MIME message ready to be sent
func (m Mail) Raw() (MIME, error) {
	var (
		err    error
		result MIME
	)

	if err = m.valid(); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	if err = m.Headers.write(writer); err != nil {
		return nil, err
	}

	if m.Parts != nil {
		for _, v := range m.Parts {
			if err = v.CreateAndWritePartTo(writer); err != nil {
				return nil, err
			}
		}
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	b := buf.Bytes()
	if bytes.Count(b, []byte("\n")) < 2 {
		return nil, Error("invalid e-mail content")
	}

	result = bytes.SplitN(b, []byte("\n"), 2)[1]

	return result, nil
}

//Sends message by SendMail function from net/smtp package
func (m *Mail) SmtpSend(addr string, auth smtp.Auth) error {
	b, err := m.Raw()
	if err != nil {
		return err
	}

	return smtp.SendMail(addr, auth, m.Headers.From.String(), m.Headers.GetRecipients(), b)
}

func (m *Mail) Subject(s string) {
	m.Headers.Subject = s
}

func (m *Mail) From(s string) error {
	return m.Headers.From.Parse(s)
}

func (m *Mail) ReplyTo(s string) error {
	return m.Headers.ReplyTo.Parse(s)
}

func (m *Mail) To(s string, list ...string) error {
	return m.Headers.To.ParseList(strings.Join(append(list, s), ","))
}

func (m *Mail) Cc(s string, list ...string) error {
	return m.Headers.Cc.ParseList(strings.Join(append(list, s), ","))
}

func (m *Mail) Bcc(s string, list ...string) error {
	return m.Headers.Bcc.ParseList(strings.Join(append(list, s), ","))
}

func (m *Mail) AppendPart(pc PartCreator) {
	m.Parts = append(m.Parts, pc)
}

func (m *Mail) SetCustomHeader(key, value string) {
	m.Headers.Set(key, value)
}

func (m *Mail) valid() error {
	switch {
	case m.Headers.From.valid() != nil:
		return Error("invalid header 'from' address")

	case len(m.Headers.To) == 0:
		return Error("invalid header 'to' address")
	}

	return nil
}

func (rm MIME) String() string {
	return string(rm)
}
