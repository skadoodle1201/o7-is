package tools

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func ParseInteger(data []byte) (Integer, []byte, error) {
	if len(data) < 4 {
		return 0, data, errors.New("integer needs at least 4 characters")
	}
	data = data[1:]
	isNegative := data[0] == '-'
	hasSign := data[0] == '+' || data[0] == '-'
	if hasSign {
		data = data[1:]
	}
	integerData := bytes.SplitN(data, []byte(CLRF), 2)
	if len(integerData) != 2 {
		return 0, data, errors.New("integer needs end delimiter")
	}
	data = integerData[1]
	integer, err := strconv.Atoi(string(integerData[0]))
	if err != nil {
		return 0, data, fmt.Errorf("value is not an integer: %s", integerData[0])
	}
	if isNegative {
		integer = -integer
	}
	return Integer(integer), data, nil
}
