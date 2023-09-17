package easymail

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
)

type Error string

func (e Error) Error() string {
	return fmt.Sprintf("easymail error: %s", string(e))
}

type MIME struct {
	Headers *Header
	Parts   []PartCreator
}

type RawMIME []byte

type PartCreator interface {
	GetHeaders() textproto.MIMEHeader
	CreateAndWritePartTo(*multipart.Writer) error
}

func NewMime() *MIME {
	return &MIME{
		Headers: &Header{
			encoder: mime.BEncoding,
			charset: "UTF-8",
			custom:  textproto.MIMEHeader{},
		},
	}
}

//Returns a MIME message ready to be sent
func (m MIME) Raw() (RawMIME, error) {
	var (
		err    error
		result RawMIME
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
func (m *MIME) SmtpSend(addr string, auth smtp.Auth) error {
	b, err := m.Raw()
	if err != nil {
		return err
	}

	return smtp.SendMail(addr, auth, m.Headers.From.String(), m.Headers.GetRecipients(), b)
}

func (m *MIME) Subject(s string) {
	m.Headers.Subject = s
}

func (m *MIME) From(s string) error {
	return m.Headers.From.Parse(s)
}

func (m *MIME) ReplyTo(s string) error {
	return m.Headers.ReplyTo.Parse(s)
}

func (m *MIME) To(s string, list ...string) error {
	return m.Headers.To.ParseList(strings.Join(append(list, s), ","))
}

func (m *MIME) Cc(s string, list ...string) error {
	return m.Headers.Cc.ParseList(strings.Join(append(list, s), ","))
}

func (m *MIME) Bcc(s string, list ...string) error {
	return m.Headers.Bcc.ParseList(strings.Join(append(list, s), ","))
}

func (m *MIME) AppendPart(pc PartCreator) {
	m.Parts = append(m.Parts, pc)
}

func (m *MIME) SetCustomHeader(key, value string) {
	m.Headers.Set(key, value)
}

func (m *MIME) valid() error {
	switch {
	case m.Headers.From.valid() != nil:
		return Error("invalid header 'from' address")

	case len(m.Headers.To) == 0:
		return Error("invalid header 'to' address")
	}

	return nil
}

func (rm RawMIME) String() string {
	return string(rm)
}
