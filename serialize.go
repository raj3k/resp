package resp

import (
	"bytes"
	"errors"
	"strconv"
)

type RespSerializer struct {
}

func NewRespSerializer() *RespSerializer {
	return &RespSerializer{}
}

func (s *RespSerializer) Serialize(value Value) ([]byte, error) {
	return value.Encode()
}

func (s *RespSerializer) Deserialize(input []byte) (Value, error) {
	if len(input) == 0 {
		return nil, errors.New("empty input")
	}

	switch input[0] {
	case '+':
		return decodeSimpleString(input[1:])
	default:
		return nil, errors.New("unknown data type")
	}
}

func decodeSimpleString(input []byte) (Value, error) {
	endIndex := bytesIndexOf(input, []byte{'\r', '\n'})
	if endIndex == -1 {
		return nil, errors.New("invalid SimpleString format")
	}
	return &SimpleString{Data: input[:endIndex]}, nil
}

// decodeError decodes an Error value from the RESP format.
func decodeError(input []byte) (*Error, error) {
	endIndex := bytesIndexOf(input, []byte{'\r', '\n'})
	if endIndex == -1 {
		return nil, errors.New("invalid Error format")
	}
	return &Error{Data: input[:endIndex]}, nil
}

// decodeInteger decodes an Integer value from the RESP format.
func decodeInteger(input []byte) (*Integer, error) {
	endIndex := bytesIndexOf(input, []byte{'\r', '\n'})
	if endIndex == -1 {
		return nil, errors.New("invalid Integer format")
	}
	return &Integer{Data: input[:endIndex]}, nil
}

// decodeBulkString decodes a BulkString value from the RESP format.
func decodeBulkString(input []byte) (*BulkString, error) {
	if len(input) < 2 {
		return nil, errors.New("invalid BulkString format")
	}
	if input[0] == '-' && input[1] == '1' && len(input) == 3 && input[2] == '\r' {
		return &BulkString{}, nil // Null BulkString
	}

	sizeEndIndex := bytesIndexOf(input, []byte{'\r', '\n'})
	if sizeEndIndex == -1 {
		return nil, errors.New("invalid BulkString format")
	}
	size, err := strconv.Atoi(string(input[:sizeEndIndex]))
	if err != nil {
		return nil, errors.New("invalid BulkString size")
	}

	dataStartIndex := sizeEndIndex + 2
	dataEndIndex := dataStartIndex + size
	if len(input) < dataEndIndex+2 || input[dataEndIndex] != '\r' || input[dataEndIndex+1] != '\n' {
		return nil, errors.New("invalid BulkString format")
	}

	return &BulkString{Data: input[dataStartIndex:dataEndIndex]}, nil
}

// decodeArray decodes an Array value from the RESP format.
func decodeArray(input []byte) (*Array, error) {
	if len(input) < 2 {
		return nil, errors.New("invalid Array format")
	}

	// Check if the array is a null array.
	if input[0] == '*' && input[1] == '-' && input[2] == '1' && len(input) == 4 && input[3] == '\r' {
		return &Array{}, nil // Null Array
	}

	// Parse the size of the array.
	sizeEndIndex := bytesIndexOf(input, []byte{'\r', '\n'})
	if sizeEndIndex == -1 {
		return nil, errors.New("invalid Array format")
	}
	size, err := strconv.Atoi(string(input[1:sizeEndIndex]))
	if err != nil {
		return nil, errors.New("invalid Array size")
	}

	serializer := NewRespSerializer()

	// Initialize variables to keep track of the array elements.
	arrayElementsStartIndex := sizeEndIndex + 2
	arrayElementsEndIndex := arrayElementsStartIndex
	var arrayElements []Value

	// Loop to decode each element in the array.
	for i := 0; i < size; i++ {
		// Decode the element using the Decode function, which can handle all data types.
		element, err := serializer.Deserialize((input[arrayElementsEndIndex:]))
		if err != nil {
			return nil, err
		}

		// Add the decoded element to the arrayElements slice.
		arrayElements = append(arrayElements, element)

		// Calculate the new endIndex for the next iteration based on the length of the encoded element.
		encodedElement, err := element.Encode()
		if err != nil {
			return nil, err
		}
		arrayElementsEndIndex += len(input[arrayElementsEndIndex:]) - len(encodedElement)
	}

	// Return the Array instance with its ArrayElements populated with the decoded values.
	return &Array{ArrayElements: arrayElements}, nil
}

// bytesIndexOf finds the index of a subarray in a byte slice.
// Returns -1 if the subarray is not found.
func bytesIndexOf(input, sub []byte) int {
	for i := 0; i <= len(input)-len(sub); i++ {
		if bytes.Equal(input[i:i+len(sub)], sub) {
			return i
		}
	}
	return -1
}
