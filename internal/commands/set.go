package commands

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	serverhelpers "github.com/skadoodle1201/o7-is/internal/serverHelpers"
	"github.com/skadoodle1201/o7-is/internal/tools"
)

type DataStore struct {
	value    string
	expiry   time.Time
	createOn time.Time
}

var SetStore = map[string]*DataStore{}
var SetStoreMux = sync.RWMutex{}

func SetCommand(args tools.Array) (string, error) {
	_, expiryTime := "", time.Time{}
	var err error
	var message string
	if len(args) != 2 && len(args) != 4 {
		fmt.Println("\n Failed at length \n", len(args))
		err = fmt.Errorf("Incorrect Input :: %v", args)
		return message, err
	}
	key, okKey := args[0].(tools.BulkString)
	value, okValue := args[1].(tools.BulkString)

	if !okKey || !okValue {
		fmt.Println("\n Failed at Key Val \n", len(args))
		err = fmt.Errorf("Incorrect Input :: %v", args[0])
		return message, err
	}
	now :=
		time.Now()
	SetStoreMux.Lock()
	if len(args) == 4 {
		typeEx, okTypeEx := args[2].(tools.BulkString)
		exVal, okExVal := args[3].(tools.BulkString)
		if !okTypeEx || !okExVal {
			fmt.Println("\n Failed at Time \n", len(args))
			err = fmt.Errorf("Incorrect Input :: %v %v", args[2], args[3])
			return message, err
		}
		timeToExpire, convertErr := strconv.Atoi(exVal.Value)
		if convertErr != nil {
			err = fmt.Errorf("Incorrect Convert Time %v %v", args[2], args[3])
			return message, err
		}
		var typeTime string = strings.ToUpper(typeEx.Value)
		switch typeTime {
		case "PX":
			{
				expiryTime = now.Add(time.Millisecond * time.Duration(timeToExpire))
			}
		case "EX":
			{
				expiryTime = now.Add(time.Second * time.Duration(timeToExpire))
			}
		}

	}

	SetStore[key.Value] = &DataStore{
		value:    value.Value,
		expiry:   expiryTime,
		createOn: now,
	}
	SetStoreMux.Unlock()

	message = tools.SimpleString("OK").Encode()
	replicaConn := tools.GetReplicaConns()
	fmt.Println("Length Of Replica", len(replicaConn))
	if len(replicaConn) > 0 && tools.GetActiverServerRole() == tools.MASTER_ROLE {
		for _, conn := range replicaConn {
			fmt.Println("In Comnn")
			serverhelpers.SendSetCommandToReplica(conn, key.Value, value.Value)
		}
	}
	return message, err
}

func GetCommand(args tools.Array) (string, error) {
	var err error
	var message string
	if len(args) > 1 {
		err = fmt.Errorf("incorrect Input :: %v", args)
		return message, err
	}
	key, okKey := args[0].(tools.BulkString)
	if !okKey {
		err = fmt.Errorf("incorrect Input :: %v", args[0])
		return message, err
	}
	SetStoreMux.RLock()
	val, ok := SetStore[key.Value]
	if !ok {
		message = fmt.Sprintf("$-1%s", tools.CLRF)
		return message, err
	}
	now := time.Now()
	fmt.Println("IS ZERO:: ", val.expiry.IsZero())
	if !val.expiry.IsZero() && now.After(val.expiry) {
		fmt.Printf("\n %v %v \n", now.After(val.expiry), val.expiry)
		delete(SetStore, key.Value)
		message = fmt.Sprintf("$-1%s", tools.CLRF)
		return message, err
	}
	SetStoreMux.RUnlock()

	message = tools.SimpleString(val.value).Encode()
	return message, err
}
