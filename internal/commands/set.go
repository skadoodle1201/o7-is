package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

type DataStore struct {
	value    string
	expiry   time.Time
	createOn time.Time
}

var SetStore = map[string]dataStore{}
var SetStoreMux = sync.RWMutex{}

func SetCommand(args tools.Array) (string, error) {
	_, expiryTime := "", time.Time{}
	var err error
	var message string
	if len(args) != 2 || len(args) != 4 {
		err = fmt.Errorf("Incorrect Input :: %v", args)
		return message, err
	}
	key, okKey := args[0].(tools.BulkString)
	value, okValue := args[1].(tools.BulkString)

	if !okKey || !okValue {
		err = fmt.Errorf("Incorrect Input :: %v", args[0])
		return message, err
	}
	now :=
		time.Now()
	SetStoreMux.Lock()
	SetStore[key.Value] = &DataStore{
		value:    value.Value,
		expiry:   now,
		createOn: now,
	}
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
