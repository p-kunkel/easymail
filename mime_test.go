package easymail

import (
	"mime"
	"testing"

	"github.com/p-kunkel/easymail/message"
	"github.com/stretchr/testify/assert"
)

func TestNewMail(t *testing.T) {
	m := New()
	if !assert.Equal(t, mime.BEncoding, m.Headers.encoder) ||
		!assert.Equal(t, "UTF-8", m.Headers.charset) {
		t.FailNow()
	}
}

func TestMarshalingMimeToRaw(t *testing.T) {
	m := New()
	m.To("ABCśŹć <test2@example.com>")
	m.From("test2@example.com")
	actual, err := m.Raw()

	if !assert.Contains(t, string(actual), "Content-Type: multipart/mixed; boundary=") ||
		!assert.Contains(t, string(actual), "From: test2@example.com") ||
		!assert.Contains(t, string(actual), "Mime-Version: 1.0") ||
		!assert.Contains(t, string(actual), "Reply-To: test2@example.com") ||
		!assert.Contains(t, string(actual), "Return-Path: test2@example.com") ||
		!assert.Contains(t, string(actual), "To: =?UTF-8?b?QUJDxZvFucSHIDx0ZXN0MkBleGFtcGxlLmNvbT4=?=") {
		t.FailNow()
	}

	if !assert.NoError(t, err) ||
		!assert.NotNil(t, actual) {
		t.FailNow()
	}
}

func TestMimeSetSubject(t *testing.T) {
	s := "test"
	m := New()

	m.Subject(s)
	if !assert.Equal(t, s, m.Headers.Subject) {
		t.FailNow()
	}
}

func TestMimeSetFrom(t *testing.T) {
	a := Address{
		Address: "test@example.com",
	}

	m := New()
	err := m.From(a.String())
	if !assert.NoError(t, err) ||
		!assert.Equal(t, a, m.Headers.From) {
		t.FailNow()
	}
}

func TestMimeSetReplyTo(t *testing.T) {
	a := Address{
		Name:    "Test",
		Address: "test@example.com",
	}

	m := New()
	err := m.ReplyTo(a.String())
	if !assert.NoError(t, err) ||
		!assert.Equal(t, a, m.Headers.ReplyTo) {
		t.FailNow()
	}
}

func TestMimeSetTo(t *testing.T) {
	a1 := Address{
		Name:    "Test",
		Address: "test@example.com",
	}

	a2 := Address{
		Name:    "Test 2",
		Address: "test-2@example.com",
	}

	m := New()
	err := m.To(a1.String(), a2.String())

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.Equal(t, 2, len(m.Headers.To), "invalid header 'to' length") &&
		!assert.NotEqual(t, m.Headers.To[0], m.Headers.To[1]) {
		t.FailNow()
	}
}

func TestMimeSetCc(t *testing.T) {
	a1 := Address{
		Name:    "Test",
		Address: "test@example.com",
	}

	a2 := Address{
		Name:    "Test 2",
		Address: "test-2@example.com",
	}

	m := New()
	err := m.Cc(a1.String(), a2.String())

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.Equal(t, 2, len(m.Headers.Cc), "invalid header 'cc' length") &&
		!assert.NotEqual(t, m.Headers.Cc[0], m.Headers.Cc[1]) {
		t.FailNow()
	}
}

func TestMimeSetBcc(t *testing.T) {
	a1 := Address{
		Name:    "Test",
		Address: "test@example.com",
	}

	a2 := Address{
		Name:    "Test 2",
		Address: "test-2@example.com",
	}

	m := New()
	err := m.Bcc(a1.String(), a2.String())

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.Equal(t, 2, len(m.Headers.Bcc), "invalid header 'bcc' length") &&
		!assert.NotEqual(t, m.Headers.Bcc[0], m.Headers.Bcc[1]) {
		t.FailNow()
	}
}

func TestMimeAppendPart(t *testing.T) {
	m := New()
	m.AppendPart(message.New("test message"))
	m.AppendPart(message.New("test message 2"))

	if !assert.Equal(t, 2, len(m.Parts), "invalid parts length") {
		t.FailNow()
	}
}

func TestMimeSetCustomHeader(t *testing.T) {
	k := "custom-header"
	v := "test-custom-header-value"
	m := New()

	m.SetCustomHeader(k, v)
	actualValue := m.Headers.custom.Get(k)

	if !assert.Equal(t, v, actualValue) {
		t.FailNow()
	}
}

func TestMimeValidation(t *testing.T) {
	m := New()

	if !assert.Error(t, m.valid()) {
		t.FailNow()
	}

	m.From("test@example.com")
	m.To("test-2@example.com")

	if !assert.NoError(t, m.valid()) {
		t.FailNow()
	}
}
