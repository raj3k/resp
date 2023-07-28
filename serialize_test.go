package resp

import (
	"testing"
)

// Test serialization for SimpleString
func TestSerializeSimpleString(t *testing.T) {
	value := &SimpleString{Data: []byte("OK")}
	expected := "+OK\r\n"

	serialized, err := Encode(value)
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
	value := &Error{Data: []byte("Error message")}
	expected := "-Error message\r\n"

	serialized, err := Encode(value)
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
	value := &Integer{Data: []byte("42")}
	expected := ":42\r\n"

	serialized, err := Encode(value)
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
	value := &BulkString{Data: []byte("Hello, World!")}
	expected := "$13\r\nHello, World!\r\n"

	serialized, err := Encode(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

func TestSerializeBulkStringNull(t *testing.T) {
	value := &BulkString{Data: nil}
	expected := "$-1\r\n"

	serialized, err := Encode(value)
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
	value := &Array{
		ArrayElements: []Value{
			&BulkString{Data: []byte("foo")},
			&BulkString{Data: []byte("bar")},
		},
	}
	expected := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"

	serialized, err := Encode(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}

func TestSerializeArrayNull(t *testing.T) {
	value := &Array{ArrayElements: nil}
	expected := "*-1\r\n"

	serialized, err := Encode(value)
	if err != nil {
		t.Errorf("Error during serialization: %v", err)
		return
	}

	if string(serialized) != expected {
		t.Errorf("Serialization failed. Expected: %q, Got: %q", expected, string(serialized))
	}
}
