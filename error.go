package easymail

import "fmt"

type Error string

func (e Error) Error() string {
	return fmt.Sprintf("easymail error: %s", string(e))
}
