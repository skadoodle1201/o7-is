package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func RedisCommands(command string, args tools.Array) string {
	command = strings.ToUpper(command)
	switch command {
	case "PING":
		{
			return fmt.Sprintf("+PONG\r\n")
		}
	case "ECHO":
		{
			message, err := EchoCommand(args)
			if err != nil {
				return fmt.Sprintf("Invalid Operation %v", command)
			}
			return message
		}
	case "SET":
		{
			message, err := SetCommand(args)
			if err != nil {
				return fmt.Sprintf("Invalid Operation %v", command)
			}
			return message
		}
	case "GET":
		{
			message, err := GetCommand(args)
			if err != nil {
				return fmt.Sprintf("Invalid Operation %v", command)
			}
			return message
		}
	default:
		{
			return fmt.Sprintf("+PONG\r\n")
		}
	}
}
