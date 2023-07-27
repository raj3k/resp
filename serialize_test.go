package resp

import (
	"fmt"
	"testing"
)

func TestSerializeRESP(t *testing.T) {
	testData := []struct {
		input    *RESPValue
		expected string
	}{
		{
			input: &RESPValue{
				Type: SimpleString,
				Data: []byte("OK"),
			},
			expected: "+OK\r\n",
		},
		{
			input: &RESPValue{
				Type: Error,
				Data: []byte("Error message"),
			},
			expected: "-Error message\r\n",
		},
		{
			input: &RESPValue{
				Type: Integer,
				Data: []byte("42"),
			},
			expected: ":42\r\n",
		},
		{
			input: &RESPValue{
				Type: BulkString,
				Data: []byte("Hello, World!"),
			},
			expected: "$13\r\nHello, World!\r\n",
		},
		{
			input: &RESPValue{
				Type: Array,
				ArrayElements: []*RESPValue{
					{Type: BulkString, Data: []byte("foo")},
					{Type: BulkString, Data: []byte("bar")},
				},
			},
			expected: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
		},
		{
			input: &RESPValue{
				Type: Array,
			},
			expected: "*-1\r\n",
		},
		{
			input: &RESPValue{
				Type: BulkString,
			},
			expected: "$-1\r\n",
		},
		{
			input: &RESPValue{
				Type: Array,
				ArrayElements: []*RESPValue{
					{
						Type: BulkString,
						Data: []byte("hello"),
					},
					{
						Type: BulkString,
					},
					{
						Type: BulkString,
						Data: []byte("world"),
					},
				},
			},
			expected: "*3\r\n$5\r\nhello\r\n$-1\r\n$5\r\nworld\r\n",
		},
	}

	for _, td := range testData {
		serializedData := SerializeRESP(td.input)
		if serializedData != td.expected {
			fmt.Printf("SerializeRESP(%q) - Test Failed:\nExpected: %q\nGot: %q\n", td.input.Type, td.expected, string(serializedData))
		} else {
			fmt.Printf("SerializeRESP(%q) - Test Passed\n", td.input.Type)
		}
	}
}
