package commands

import (
	"fmt"

	"github.com/skadoodle1201/o7-is/internal/tools"
)

func InfoCommand(args tools.Array, role string) (message string, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("Invalid Input %v", args)
		return message, err
	}
	_, okVal := args[0].(tools.BulkString)

	fmt.Println(role)

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
master_replid:%s
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:`, roleVal, tools.ServerUUID())

	message = message + tools.RedisBulkString(valuesToRespond).Encode()

	return message, err
}
