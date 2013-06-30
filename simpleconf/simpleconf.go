// Library package for reading configuration files
package simpleconf

import (
	"fmt"
	"io"
	"os"
)

const max_key_size int = 256
const max_value_size int = 1024

type BufferParserState int

const (
	// New line received or new file
	BPS_New BufferParserState = iota
	// Reading key
	BPS_Key
	// Reading value
	BPS_Value
	// Comment
	BPS_Comment
)

type SimpleConfSettings struct {
	settings map[string] string
}

func NewSimpleConfSettings () *SimpleConfSettings {
	s := SimpleConfSettings{}
	s.settings = make(map[string] string)
	return &s
}

func (s *SimpleConfSettings) Get (key string, default_value string) string {
	if (*s).settings[key] == "" {
		return default_value
	}
	return (*s).settings[key]
}

func ParseFile (f string) (*SimpleConfSettings, error) {
	s := NewSimpleConfSettings()

	fi, err := os.Open(f)
	if err != nil {
		return s, err
	}

	defer fi.Close()

	buf := make([]byte, 1024)
	bufParser, bufCloser := MakeParseBuffer()

	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			return s, nil
		}
		if n == 0 {
			break
		}
		
		_, err = bufParser(buf, n, s)
		
		if err != nil {
			return s, err
		}
	}
	_, err = bufCloser(s)

	return s, err
}

type SimpleConfError struct {
	Where int
	What string
}

func (e SimpleConfError) Error() string {
	return fmt.Sprintf("syntax error on line %d: %s", e.Where, e.What)
}

func MakeParseBuffer () (func ([]byte, int, *SimpleConfSettings) (int, error), func (*SimpleConfSettings) (int, error)) {
	key := make([]byte, 0, max_key_size)
	value := make([]byte, 0, max_value_size)
	state := BPS_New
	line := 1
	
	return func (b []byte, n int, s *SimpleConfSettings) (int, error) {
		i := 0
		for ; i < n; i++ {
			if b[i] == 0 { // null
				break
			} else if b[i] == 10 || b[i] == 13 { // \n or \r
				if state == BPS_Key {
					return i, SimpleConfError{line, "expected equals sign before end of line"}
				} else if state == BPS_Value {
					(*s).settings[string(key)] = string(value)
				}
				key = key[:0]
				value = value[:0]
				state = BPS_New
				line += 1
			} else if state == BPS_Comment {
				continue
			} else if b[i] == 32 || b[i] == 9 { // space or tab
				if state == BPS_Value {
					if len(value) == cap(value) {
						return i, SimpleConfError{line, fmt.Sprintf("value is too long (max: %d bytes)", max_value_size)}
					}
					value = append(value, b[i])
				} else {
					return i, SimpleConfError{line, "whitespace not allowed here"}
				}
			} else if b[i] == 59 || b[i] == 91 { // semicolon, left square bracket
				state = BPS_Comment
			} else if b[i] == 61 { // equals
				if state == BPS_Key {
					state = BPS_Value
				} else {
					return i, SimpleConfError{line, "unexpected equals sign"}
				}
			} else {
				if state == BPS_New {
					state = BPS_Key
				}
				if state == BPS_Key {
					if len(key) == cap(key) {
						return i, SimpleConfError{line, fmt.Sprintf("key is too long (max: %d bytes)", max_key_size)}
					}
					key = append(key, b[i])
				} else if state == BPS_Value {
					if len(value) == cap(value) {
						return i, SimpleConfError{line, fmt.Sprintf("value is too long (max: %d bytes)", max_value_size)}
					}
					value = append(value, b[i])
				}
			}
		}
		return i, nil
	}, func (s *SimpleConfSettings) (int, error) {
		if state == BPS_Key {
			return 0, SimpleConfError{line, "unexpected equals sign"}
		} else if state == BPS_Value {
			(*s).settings[string(key)] = string(value)
		}
		return (1 + len(key) + len(value)), nil
	}

}
