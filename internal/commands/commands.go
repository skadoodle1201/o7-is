package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func RedisCommands(command string, args tools.Array, role string) string {
	command = strings.ToUpper(command)
	switch command {
	case "PING":
		{
			return fmt.Sprintf("+PONG%s", tools.CLRF)
		}
	case "ECHO":
		{
			message, err := EchoCommand(args)
			if err != nil {
				return fmt.Sprintf("-ERR Invalid Operation %v%s", command, tools.CLRF)
			}
			return message
		}
	case "SET":
		{
			message, err := SetCommand(args)
			if err != nil {
				return fmt.Sprintf("-ERR Invalid Operation %v%s", command, tools.CLRF)
			}
			return message
		}
	case "GET":
		{
			message, err := GetCommand(args)
			if err != nil {
				return fmt.Sprintf("-ERR Invalid Operation %v%s", command, tools.CLRF)
			}
			return message
		}
	case "INFO":
		{
			message, err := InfoCommand(args, role)
			if err != nil {
				return fmt.Sprintf("-ERR Invalid Operation %v%s", command, tools.CLRF)
			}
			return message

		}
	default:
		{
			return fmt.Sprintf("-ERR Invalid Operation %v%s", command, tools.CLRF)
		}
	}
}
