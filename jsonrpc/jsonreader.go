package jsonrpc

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
	done   chan []byte
}

// NewJSONReader creates a new JSONReader instance
func NewJSONReader(done chan []byte) JSONReader {
	return JSONReader{done: done}
}

func (r *JSONReader) reset() {
	r.state = state0
	r.stack = nil
	r.buffer = nil

	// log.Println("jsonreader: reset")
}

func (r *JSONReader) send() {
	if r.done != nil {
		r.done <- r.buffer
	}

	r.reset()
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
	}
}

// FeedByte feeds the JSONReader a single byte
func (r *JSONReader) FeedByte(b byte) {
	r.buffer = append(r.buffer, b)
	r.transition(b)
}

// FeedBytes feeds the JSONReader a slice of bytes
func (r *JSONReader) FeedBytes(bs []byte) {
	for _, b := range bs {
		r.FeedByte(b)
	}
}
