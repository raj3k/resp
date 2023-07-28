package resp

import (
	"fmt"
)

func Encode(value Value) ([]byte, error) {
	switch value.GetType() {
	case SimpleStringType, ErrorType, IntegerType, BulkStringType:
		return value.Encode()
	case ArrayType:
		return value.Encode()
	default:
		return nil, fmt.Errorf("Unsupported data type")
	}
}
