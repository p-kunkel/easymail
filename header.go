package gosimplemime

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
)

type Header struct {
	Subject string

	From    Address
	ReplyTo Address
	To      Addresses
	Cc      Addresses
	Bcc     Addresses

	custom  textproto.MIMEHeader
	charset string
	encoder HeaderEncoder
}

type HeaderEncoder interface {
	Encode(charset string, s string) string
}

func (h *Header) GetRecipients() []string {
	return append(append(h.To.GetListOfAddresses(), h.Cc.GetListOfAddresses()...), h.Bcc.GetListOfAddresses()...)
}

func (h *Header) Set(key, value string) {
	h.custom.Set(key, h.encoder.Encode(h.charset, value))
}

func (h *Header) SetEncoder(he HeaderEncoder, charset string) {
	h.charset = charset
	h.encoder = he
}

func (h *Header) write(mpw *multipart.Writer) error {
	if h.custom == nil {
		h.custom = make(textproto.MIMEHeader)
	}

	s := ""
	h.custom.Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=\"%s\"", mpw.Boundary()))
	h.custom.Set("MIME-Version", "1.0")

	if s = h.From.String(); s != "" {
		h.custom.Set("From", h.encoder.Encode(h.charset, s))
		h.custom.Set("Reply-To", h.encoder.Encode(h.charset, s))
	}

	if s = h.To.String(); s != "" {
		h.custom.Set("To", h.encoder.Encode(h.charset, s))
	}

	if s = h.Cc.String(); s != "" {
		h.custom.Set("Cc", h.encoder.Encode(h.charset, s))
	}

	if s = h.Bcc.String(); s != "" {
		h.custom.Set("Bcc", h.encoder.Encode(h.charset, s))
	}

	if s = h.From.String(); s != "" {
		h.custom.Set("Return-Path", h.encoder.Encode(h.charset, s))
	}

	if s = h.Subject; s != "" {
		h.custom.Set("Subject", h.encoder.Encode(h.charset, s))
	}

	if h.ReplyTo.valid() == nil {
		h.custom.Set("Reply-To", h.encoder.Encode(h.charset, h.ReplyTo.String()))
	}

	_, err := mpw.CreatePart(h.custom)
	return err
}
