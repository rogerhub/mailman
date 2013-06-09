/**
 *  Unit tests for ExpandingBuffer
 */

package mailman

import (
	"testing"
)

func TestNewExpandingBuffer (t *testing.T) {
	e := NewExpandingBuffer()
	if e == nil {
		t.Error("new buffer creation failed")
	}
	if e.Length != 0 {
		t.Errorf("new buffer does not have length 0 (length is %d)", e.Length)
	}
}

func TestExpandingBuffer (t *testing.T) {
	e := NewExpandingBuffer()
	r := e.InsertItem(4)
	if 0 != r {
		t.Errorf("added 1 element to new list: return value should be 0 (return was %d)", r)
	}
	i := e.GetItem(0)
	if i != 4 {
		t.Errorf("tried to get 1st element of 1-element list, but got incorrect value: expecting 4, got: %d", i)
	}
	for i := 1; i <= 100; i ++ {
		e.InsertItem(i + 4)
	}
	if e.Length != 101 {
		t.Errorf("length of 101-element list should be 101 (length was %d)", e.Length)
	}
	for i := 0; i < e.Length; i ++ {
		if e.GetItem(i) != i + 4 {
			t.Errorf("tried to insert 100 elements into list, but element at index %d does not have value %d (returned was %d)", i, i + 4, e.GetItem(i))
		}
	}
}
