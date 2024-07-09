package tools

import "fmt"

func (s SimpleString) Encode() string {
	return fmt.Sprintf("+%s%s", s, CLRF)
}
