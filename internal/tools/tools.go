package tools

import (
	"fmt"
)

func Parse(data []byte) (interface{}, []byte, error) {
	dataType := data[0]
	var (
		parsingType string
		parsed      interface{}
		err         error
	)
	switch dataType {
	case '*':
		parsingType = "Array"
		parsed, data, err = ParseArray(data)
	case '$':
		parsingType = "BulkString"
		parsed, data, err = ParseBulkString(data)
	case ':':
		parsingType = "Integer"
		parsed, data, err = ParseInteger(data)
	case '+':
		return nil, data, nil
	default:
		return nil, data, fmt.Errorf("unknown data type %v", string(dataType))
	}
	if err != nil {
		return nil, data, fmt.Errorf("parsing type %s: %w", parsingType, err)
	}
	return parsed, data, nil
}
