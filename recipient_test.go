package easymail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipientsListParsing(t *testing.T) {
	var (
		name1    = "Jonh"
		email1   = "john@example.com"
		address1 = fmt.Sprintf("%s <%s>", name1, email1)

		name2    = "Alice"
		email2   = "alice@example.com"
		address2 = fmt.Sprintf("%s <%s>", name2, email2)

		r Addresses
	)

	err := r.ParseList(fmt.Sprintf("%s, %s", address1, address2))

	if !assert.NoError(t, err) ||
		!assert.Equal(t, r[0].Name, name1) ||
		!assert.Equal(t, r[0].Address, email1) {

		t.FailNow()
	}
}

func TestRecipientsAppend(t *testing.T) {
	var (
		name1    = "Jonh"
		email1   = "john@example.com"
		address1 = fmt.Sprintf("%s <%s>", name1, email1)
		r        Addresses
	)

	err := r.Append(address1)

	if !assert.NoError(t, err) ||
		!assert.Equal(t, r[0].Name, name1) ||
		!assert.Equal(t, r[0].Address, email1) {

		t.FailNow()
	}
}
