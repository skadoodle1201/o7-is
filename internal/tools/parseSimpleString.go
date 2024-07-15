package tools

import "fmt"

func (s SimpleString) Encode() string {
	return fmt.Sprintf("+%s%s", s, CLRF)
}

func (s RedisBulkString) Encode() string {
	lengthOfString := len(s)
	return fmt.Sprintf("$%d%s%s%s", lengthOfString, CLRF, s, CLRF)
}
