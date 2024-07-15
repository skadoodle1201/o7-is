package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func ReplConfCommand(args tools.Array) (message string, err error) {
	message = tools.SimpleString("OK").Encode()
	return message, err
}

func PsyncConfCommand(args tools.Array) (message string, err error) {
	message = tools.SimpleString(fmt.Sprintf("FULLRESYNC %s 0", tools.ServerUUID())).Encode()
	return message, err
}
