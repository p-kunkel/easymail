package gosimplemime

import (
	"io"
	"mime/quotedprintable"
)

func writeQuotedPrintable(w io.Writer, data []byte) error {
	quotedPrintableEncoder := quotedprintable.NewWriter(w)
	defer quotedPrintableEncoder.Close()
	_, err := quotedPrintableEncoder.Write(data)
	return err
}
