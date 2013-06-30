// Test cases for simple configuration package
package simpleconf

import (
	"testing"
)

func (s *SimpleConfSettings) AssertIs (t *testing.T, key string, value string) {
	if s.Get(key, "not " + value) != value {
		t.Errorf("Failed assertion: value of %s is %s (expected %s)", key, s.Get(key, "not " + value), value)
	}
}

func SetupTestBuffer (t *testing.T, data string) *SimpleConfSettings {
	bufParser, bufCloser := MakeParseBuffer()
	buf := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		buf[i] = data[i]
	}
	s := NewSimpleConfSettings()

	_, err := bufParser(buf, len(buf), s)
	if err != nil {
		t.Errorf("error raised %s", err.Error())
	}
	_, err = bufCloser(s)
	if err != nil {
		t.Errorf("error raised %s", err.Error())
	}

	return s
}

func AssertSyntaxError (t *testing.T, data string) {
	bufParser, bufCloser := MakeParseBuffer()
	buf := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		buf[i] = data[i]
	}
	s := NewSimpleConfSettings()
	
	_, err1 := bufParser(buf, len(buf), s)
	_, err2 := bufCloser(s)
	if err1 == nil && err2 == nil {
		t.Errorf("expected syntax error but got none (source follows)\n\n%s\n\n", data)
	}
}

func TestBufferParser (t *testing.T) {
	test1 := SetupTestBuffer(t, "one=1\ntwo=2\nthree=3")
	test1.AssertIs(t, "one", "1")
	test1.AssertIs(t, "three", "3")

	test2 := SetupTestBuffer(t, "; comment\n; comment\r\n[header]\n\nOne=1 uno\n\nTwo=2 dos\n\nThree=3 tres\n\n; comment")
	test2.AssertIs(t, "One", "1 uno")
	test2.AssertIs(t, "Two", "2 dos")
	test2.AssertIs(t, "Three", "3 tres")
	
	AssertSyntaxError(t, "=value")
	AssertSyntaxError(t, "key")
	AssertSyntaxError(t, "no space in key=value")
	AssertSyntaxError(t, "key=value=value")
}
