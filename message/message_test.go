package message

import (
	"bytes"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageWrite(t *testing.T) {
	b := new(bytes.Buffer)
	mw := multipart.NewWriter(b)
	m := New("msg śćź")

	err := m.CreateAndWritePartTo(mw)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.Contains(t, b.String(), "msg =C5=9B=C4=87=C5=BA") {
		t.FailNow()
	}
}

func TestMessageGetHeaders(t *testing.T) {
	m := New("test")
	h := m.GetHeaders()

	if !assert.Equal(t, "text/plain; charset=utf-8", h.Get("content-type")) ||
		!assert.Equal(t, "quoted-printable", h.Get("Content-Transfer-Encoding")) {
		t.FailNow()
	}

	m = New("<!DOCTYPE html> test </html>")
	h = m.GetHeaders()

	if !assert.Equal(t, "text/html; charset=utf-8", h.Get("content-type")) ||
		!assert.Equal(t, "quoted-printable", h.Get("Content-Transfer-Encoding")) {
		t.FailNow()
	}
}
