package tools

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func ParseBulkString(data []byte) (BulkString, []byte, error) {
	bulkString := BulkString{
		Value:  "",
		IsNull: false,
	}
	if len(data) < 6 {
		return bulkString, data, errors.New("bulk string needs at least 6 chracter")
	}
	data = data[1:]
	bulkStrData := bytes.SplitN(data, []byte(CLRF), 2)

	if len(bulkStrData) != 2 {
		return bulkString, data, errors.New("bulk string needs start delimiter")
	}
	data = bulkStrData[1]
	length, err := strconv.Atoi(string(bulkStrData[0]))
	if err != nil {
		return bulkString, data, fmt.Errorf("length is not an integer: %s", bulkStrData[0])
	}
	if length == -1 {
		bulkString.IsNull = true
		return bulkString, data, nil
	}
	fmt.Println("data: IN BULK STRING", string(data))
	bulkStrData = bytes.SplitN(data, []byte(CLRF), 2)
	if len(bulkStrData) != 2 {
		return bulkString, data, errors.New("bulk string needs end delimiter")
	}
	data = bulkStrData[1]
	bulkString.Value = string(bulkStrData[0])
	return bulkString, data, nil
}
func (d BulkString) Encode() string {
	if d.IsNull {
		return "$-1" + CLRF
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(d.Value), d.Value)
}
