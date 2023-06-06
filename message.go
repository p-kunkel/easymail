package gosimplemime

import (
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/textproto"
)

type Message []byte

func NewMessage(s string) *Message {
	result := Message(s)
	return &result
}

func (m *Message) write(mpw *multipart.Writer) error {
	if m == nil {
		m = NewMessage("")
	}

	cp, err := mpw.CreatePart(m.getHeaders())
	if err != nil {
		return err
	}

	return writeQuotedPrintable(cp, *m)
}

func (m *Message) getHeaders() textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	h.Set("Content-Type", http.DetectContentType(*m))
	return h
}

func writeQuotedPrintable(w io.Writer, data []byte) error {
	quotedPrintableEncoder := quotedprintable.NewWriter(w)
	defer quotedPrintableEncoder.Close()
	_, err := quotedPrintableEncoder.Write(data)
	return err
}
