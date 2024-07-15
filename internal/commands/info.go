package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func InfoCommand(args tools.Array, role string) (message string, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("Invalid Input %v", args)
		return message, err
	}
	_, okVal := args[0].(tools.BulkString)

	if !okVal {
		err = fmt.Errorf("Invalid Input %v", args)
		return message, err
	}

	roleVal := "master"
	if role != "master" {
		roleVal = "slave"
	}
	valuesToRespond := fmt.Sprintf(`# Replication
role:%s
connected_slaves:0
master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:`, roleVal)

	message = message + tools.RedisBulkString(valuesToRespond).Encode()

	return message, err
}
