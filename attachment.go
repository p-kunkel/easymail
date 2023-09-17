package easymail

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
)

type Attachment struct {
	ContentType string
	Filename    string
	File        bytes.Buffer
}

func NewAttachment() *Attachment {
	return &Attachment{
		File: bytes.Buffer{},
	}
}

func (a *Attachment) ReadFile(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if _, err = a.Write(f); err != nil {
		return err
	}

	if a.Filename == "" {
		a.Filename = filepath.Base(path)
	}

	if a.ContentType == "" {
		a.DetectContentType()
	}
	return nil
}

func (a *Attachment) Write(b []byte) (int, error) {
	return a.File.Write(b)
}

func (a *Attachment) DetectContentType() {
	a.ContentType = http.DetectContentType(a.File.Bytes())
}

func (a *Attachment) CreateAndWritePartTo(mpw *multipart.Writer) error {
	cp, err := mpw.CreatePart(a.GetHeaders())
	if err != nil {
		return err
	}

	dst := make([]byte, base64.StdEncoding.EncodedLen(a.File.Len()))
	base64.StdEncoding.Encode(dst, a.File.Bytes())
	_, err = cp.Write(dst)
	return err
}

func (a *Attachment) GetHeaders() textproto.MIMEHeader {
	if a.ContentType == "" {
		a.DetectContentType()
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "attachment; filename=\""+a.Filename+"\"")
	h.Set("Content-Type", a.ContentType+"; name=\""+a.Filename+"\"")
	h.Set("Content-Transfer-Encoding", "base64")
	return h
}
