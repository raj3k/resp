package resp

import (
	"errors"
	"testing"
)

// Test serialization for SimpleString
func TestSerializeSimpleString(t *testing.T) {
	serializer := NewRespSerializer()

	value := &SimpleString{Data: []byte("OK")}
	expected := "+OK\r\n"

	serialized, err := serializer.Serialize(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

// Test serialization for Error
func TestSerializeError(t *testing.T) {
	serializer := NewRespSerializer()

	value := &Error{Data: []byte("Error message")}
	expected := "-Error message\r\n"

	serialized, err := serializer.Serialize(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

// Test serialization for Integer
func TestSerializeInteger(t *testing.T) {
	serializer := NewRespSerializer()

	value := &Integer{Data: []byte("42")}
	expected := ":42\r\n"

	serialized, err := serializer.Serialize(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

// Test serialization for BulkString
func TestSerializeBulkString(t *testing.T) {
	serializer := NewRespSerializer()

	value := &BulkString{Data: []byte("Hello, World!")}
	expected := "$13\r\nHello, World!\r\n"

	serialized, err := serializer.Serialize(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

func TestSerializeBulkStringNull(t *testing.T) {
	serializer := NewRespSerializer()

	value := &BulkString{Data: nil}
	expected := "$-1\r\n"

	serialized, err := serializer.Serialize(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

// Test serialization for Array
func TestSerializeArray(t *testing.T) {
	serializer := NewRespSerializer()

	value := &Array{
		ArrayElements: []Value{
			&BulkString{Data: []byte("foo")},
			&BulkString{Data: []byte("bar")},
		},
	}
	expected := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"

	serialized, err := serializer.Serialize(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

func TestSerializeArrayNull(t *testing.T) {
	serializer := NewRespSerializer()

	value := &Array{ArrayElements: nil}
	expected := "*-1\r\n"

	serialized, err := serializer.Serialize(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

func TestDecodeSimpleString(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected Value
		err      error
	}{
		// Valid input
		{
			input:    []byte("+Hello, World!\r\n"),
			expected: &SimpleString{Data: []byte("Hello, World!")},
			err:      nil,
		},

		// Empty input
		{
			input:    []byte(""),
			expected: nil,
			err:      errors.New("invalid SimpleString format"),
		},

		// Missing newline
		{
			input:    []byte("+Hello, World!\r"),
			expected: nil,
			err:      errors.New("invalid SimpleString format"),
		},

		// Missing "+"
		{
			input:    []byte("Hello, World!\r\n"),
			expected: nil,
			err:      errors.New("invalid SimpleString format"),
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.input), func(t *testing.T) {
			value, err := decodeSimpleString(tc.input)
			if tc.err != nil {
				if err == nil || err.Error() != tc.err.Error() {
					t.Errorf("Expected error: %v, got: %v", tc.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if value.GetType() != tc.expected.GetType() {
					t.Errorf("Expected type: %v, got: %v", tc.expected.GetType(), value.GetType())
				}
				// Compare the SimpleString Data field.
				expectedData := tc.expected.(*SimpleString).Data
				actualData := value.(*SimpleString).Data
				if string(expectedData) != string(actualData) {
					t.Errorf("Expected Data: %q, got: %q", expectedData, actualData)
				}
			}
		})
	}
}
