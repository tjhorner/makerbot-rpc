package jsonrpc

import (
	"sync"
)

type jsonReaderState int

const (
	state0 jsonReaderState = iota
	state1
	state2
	state3
	state4
)

// JSONReader is a Go re-implementation of the JsonReader from
// MakerBot's C++ jsonrpc library:
// https://github.com/makerbot/jsonrpc/blob/develop/src/main/cpp/jsonreader.cpp
type JSONReader struct {
	state  jsonReaderState
	stack  []byte
	buffer []byte
	done   func([]byte) error
	rawCh  *chan []byte
	rawExp int
	mux    sync.Mutex
}

// NewJSONReader creates a new JSONReader instance
func NewJSONReader(done func([]byte) error) JSONReader {
	return JSONReader{done: done}
}

func (r *JSONReader) reset() {
	r.state = state0
	r.stack = nil
	r.buffer = nil
	r.rawCh = nil
	r.rawExp = 0
}

func (r *JSONReader) send() {
	if r.rawCh != nil || r.state == state4 {
		return
	}

	err := r.done(r.buffer)
	if err == nil {
		r.reset()
	}
}

func (r *JSONReader) transition(b byte) {
	switch r.state {
	case state0:
		if b == '{' || b == '[' {
			r.state = state1
			r.stack = append(r.stack, b)
		} else if b != ' ' && b != '\t' && b != '\n' && b != '\r' {
			r.send()
		}
		break

	case state1:
		if b == '"' {
			r.state = state2
		} else if b == '{' || b == '[' {
			r.stack = append(r.stack, b)
		} else if b == '}' || b == ']' {
			send := false

			if len(r.stack) == 0 {
				send = true
			} else {
				fch := r.stack[len(r.stack)-1]
				r.stack = r.stack[:len(r.stack)-1]

				if (fch == '{' && b != '}') || (fch == '[' && b != ']') {
					send = true
				} else {
					send = len(r.stack) == 0
				}
			}

			if send {
				r.send()
			}
		}
		break

	case state2:
		if b == '"' {
			r.state = state1
		} else if b == '\\' {
			r.state = state3
		}
		break

	case state3:
		r.state = state2
		break

	case state4:
		if len(r.buffer) >= r.rawExp {
			if r.rawCh != nil {
				*r.rawCh <- r.buffer
			}

			r.reset()
		}
		break
	}
}

// Reset resets
func (r *JSONReader) Reset() {
	r.reset()
}

// feedByte feeds the JSONReader a single byte
func (r *JSONReader) feedByte(b byte) {
	// FIXME: WHY IS SO MUCH MEMORY BEING ALLOCATED HERE???? it makes me so sad
	// https://user-images.githubusercontent.com/2646487/57392551-f1999f80-7175-11e9-9fd7-1da09bf334dd.png
	if r != nil { // FIXME: we get segfaults without this check... which shouldn't happen
		r.mux.Lock()
		defer r.mux.Unlock()
		r.transition(b)
	}
}

// Write feeds the JSONReader a slice of bytes
func (r *JSONReader) Write(bs []byte) (n int, err error) {
	if r != nil {
		r.buffer = append(r.buffer, bs...)

		for _, b := range bs {
			r.feedByte(b)
		}
	}

	return len(bs), nil
}

// GetRawData grabs raw data from the TCP connection until
// `length` is reached. The captured data is returned as an
// array of bytes.
func (r *JSONReader) GetRawData(length int) []byte {
	// prevent anything from writing while we set up for raw reading
	r.mux.Lock()

	ch := make(chan []byte)
	r.rawCh = &ch
	r.state = state4
	r.rawExp = length

	r.mux.Unlock()

	data := <-ch
	close(ch)

	return data
}
