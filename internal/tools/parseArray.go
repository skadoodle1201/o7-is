package tools

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func ParseArray(data []byte) (Array, []byte, error) {
	if len(data) < 4 {
		return nil, data, errors.New("array needs at least 4 chracter")
	}
	data = data[1:]
	arrayData := bytes.SplitN(data, []byte(CLRF), 2)
	if len(arrayData) != 2 {
		return nil, data, errors.New("array needs start delimiter")
	}
	data = arrayData[1]
	length, err := strconv.Atoi(string(arrayData[0]))
	if err != nil {
		return nil, data, fmt.Errorf("length is not an integer: %s", arrayData[0])
	}
	array := Array{}
	for i := 0; i < length; i++ {
		arrayItem, dataLeft, err := Parse(data)
		if err != nil {
			return nil, data, fmt.Errorf("failed to parse array item: %s, due to %w", data, err)
		}
		array = append(array, arrayItem)
		data = dataLeft
	}
	return array, data, nil
}
