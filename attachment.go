package gosimplemime

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type Attachment struct {
	ContentType string
	Filename    string
	File        bytes.Buffer
}

func (a *Attachment) CreateAndWritePartTo(mpw *multipart.Writer) error {
	cp, err := mpw.CreatePart(a.GetHeaders())
	if err != nil {
		return err
	}

	dst := []byte{}
	base64.StdEncoding.Encode(dst, a.File.Bytes())
	_, err = cp.Write(dst)
	return err
}

func (a *Attachment) GetHeaders() textproto.MIMEHeader {
	if a.ContentType == "" {
		a.ContentType = http.DetectContentType(a.File.Bytes())
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "attachment; filename=\""+a.Filename+"\"")
	h.Set("Content-Type", a.ContentType+"; name=\""+a.Filename+"\"")
	h.Set("Content-Transfer-Encoding", "base64")
	return h
}
