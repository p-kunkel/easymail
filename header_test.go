package easymail

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderGetRecipients(t *testing.T) {
	var (
		h Header

		emailTo  = "to@test.com"
		emailCc  = "cc@test.com"
		emailBcc = "bcc@test.com"
	)

	h.To.Append(emailTo)
	h.Cc.Append(emailCc)
	h.Bcc.Append(emailBcc)

	recipients := h.GetRecipients()
	expectedEmails := map[string]bool{
		emailTo:  false,
		emailCc:  false,
		emailBcc: false,
	}

	for _, v := range recipients {
		expectedEmails[v] = true
	}

	for email, emailExist := range expectedEmails {
		if !assert.True(t, emailExist, fmt.Sprintf("not found '%s' email in recippients", email)) {
			t.FailNow()
		}
	}
}

func TestHeaderWrite(t *testing.T) {
	var (
		h   Header
		b   = bytes.Buffer{}
		mpw = multipart.NewWriter(&b)
	)
	mpw.SetBoundary("test123")

	h.write(mpw)
	s := b.String()
	if !assert.Contains(t, s, "--test123") ||
		!assert.Contains(t, s, "Content-Type: multipart/mixed; boundary=\"test123\"") ||
		!assert.Contains(t, s, "Mime-Version: 1.0") {
		t.FailNow()
	}
}
