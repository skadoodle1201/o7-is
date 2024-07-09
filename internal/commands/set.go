package commands

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

var SetStore = map[string]string{}
var SetStoreMux = sync.RWMutex{}

func SetCommand(args tools.Array) (string, error) {
	var err error
	var message string
	if len(args) != 2 {
		err = fmt.Errorf("Incorrect Input :: %v", args)
		return message, err
	}
	key, okKey := args[0].(tools.BulkString)
	value, okValue := args[1].(tools.BulkString)
	if !okKey || !okValue {
		err = fmt.Errorf("Incorrect Input :: %v", args[0])
		return message, err
	}
	SetStoreMux.Lock()
	SetStore[key.Value] = value.Value
	SetStoreMux.Unlock()

	message = tools.SimpleString("OK").Encode()
	return message, err
}

func GetCommand(args tools.Array) (string, error) {
	var err error
	var message string
	if len(args) > 1 {
		err = fmt.Errorf("Incorrect Input :: %v", args)
		return message, err
	}
	key, okKey := args[0].(tools.BulkString)
	if !okKey {
		err = fmt.Errorf("Incorrect Input :: %v", args[0])
		return message, err
	}
	SetStoreMux.RLock()
	val, ok := SetStore[key.Value]
	SetStoreMux.RUnlock()
	if !ok {
		message = tools.SimpleString("$-1\r\n").Encode()
		return message, err
	}

	message = tools.SimpleString(val).Encode()
	return message, err
}
