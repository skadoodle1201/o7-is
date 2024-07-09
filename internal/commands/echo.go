package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func EchoCommand(args tools.Array) (string, error) {
	var err error
	var message string
	if len(args) != 1 {
		err = fmt.Errorf("Incorrect Input :: %v", args)
		return message, err
	}
	echoed, ok := args[0].(tools.BulkString)
	if !ok {
		err = fmt.Errorf("Incorrect Input :: %v", args[0])
		return message, err
	}

	message = tools.SimpleString(echoed.Value).Encode()
	return message, err
}
