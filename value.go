package resp

import (
	"bytes"
	"fmt"
)

// DataType represents the possible data types in the RESP format.
type DataType string

const (
	SimpleStringType DataType = "SimpleString"
	ErrorType        DataType = "Error"
	IntegerType      DataType = "Integer"
	BulkStringType   DataType = "BulkString"
	ArrayType        DataType = "Array"
)

// Value is an interface representing a value in the RESP format.
type Value interface {
	GetType() DataType       // GetType returns the data type of the value.
	Encode() ([]byte, error) // Encode encodes the value to RESP format.
}

// SimpleString represents a simple string value in RESP format.
type SimpleString struct {
	Data []byte
}

func (v *SimpleString) GetType() DataType {
	return SimpleStringType
}

func (v *SimpleString) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf("+%s\r\n", v.Data)), nil
}

// Error represents an error value in RESP format.
type Error struct {
	Data []byte
}

func (v *Error) GetType() DataType {
	return ErrorType
}

func (v *Error) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf("-%s\r\n", v.Data)), nil
}

// Integer represents an integer value in RESP format.
type Integer struct {
	Data []byte
}

func (v *Integer) GetType() DataType {
	return IntegerType
}

func (v *Integer) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf(":%s\r\n", v.Data)), nil
}

// BulkString represents a bulk string value in RESP format.
type BulkString struct {
	Data []byte
}

func (v *BulkString) GetType() DataType {
	return BulkStringType
}

func (v *BulkString) IsNull() bool {
	return v.Data == nil
}

func (v *BulkString) Encode() ([]byte, error) {
	if v.IsNull() {
		return []byte("$-1\r\n"), nil // Null BulkString
	}
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Data), v.Data)), nil
}

// Array represents an array value in RESP format.
type Array struct {
	ArrayElements []Value
}

func (v *Array) GetType() DataType {
	return ArrayType
}

func (v *Array) IsNull() bool {
	return v.ArrayElements == nil
}

func (v *Array) Encode() ([]byte, error) {
	if v.IsNull() {
		return []byte("*-1\r\n"), nil // Null Array
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("*%d\r\n", len(v.ArrayElements)))
	for _, element := range v.ArrayElements {
		serialized, err := element.Encode()
		if err != nil {
			return nil, err
		}
		buffer.Write(serialized)
	}
	return buffer.Bytes(), nil
}
