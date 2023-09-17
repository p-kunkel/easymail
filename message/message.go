package message

import (
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/textproto"
)

type Message []byte

func New(s string) *Message {
	result := Message(s)
	return &result
}

func (m *Message) CreateAndWritePartTo(mpw *multipart.Writer) error {
	if m == nil {
		m = New("")
	}

	partWriter, err := mpw.CreatePart(m.GetHeaders())
	if err != nil {
		return err
	}

	qp := quotedprintable.NewWriter(partWriter)
	defer qp.Close()
	_, err = qp.Write(*m)
	return err
}

func (m *Message) GetHeaders() textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	h.Set("Content-Type", http.DetectContentType(*m))
	return h
}
