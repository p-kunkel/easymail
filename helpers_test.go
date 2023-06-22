package gosimplemime

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteQuotedPrintable(t *testing.T) {
	b := new(bytes.Buffer)
	data := []byte("test ĄźĆ=")

	err := writeQuotedPrintable(b, data)
	if !assert.NoError(t, err) ||
		!assert.Equal(t, "test =C4=84=C5=BA=C4=86=3D", b.String()) {
		t.FailNow()
	}
}
