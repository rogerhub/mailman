// List class designed to support any-sized lists but minimize memory usage with 1-length list (most common).
package server

type ExpandingBufferNode interface {}
type ExpandingBuffer struct {
	Length int
	Contents []ExpandingBufferNode
}

func NewExpandingBuffer () (*ExpandingBuffer) {
	e := new(ExpandingBuffer)
	e.Contents = make([]ExpandingBufferNode, 1)
	return e
}

func (e *ExpandingBuffer) GetItem (index int) (ExpandingBufferNode) {
	if index > e.Length {
		return nil
	}
	return e.Contents[index]
}

func (e *ExpandingBuffer) InsertItem (n ExpandingBufferNode) (int) {
	if e.Length == len(e.Contents) {
		c := make([]ExpandingBufferNode, 2 * e.Length)
		for i := 0; i < e.Length; i ++ {
			c[i] = e.Contents[i]
		}
		e.Contents = c
	}
	e.Contents[e.Length] = n
	e.Length ++
	return e.Length - 1
}
