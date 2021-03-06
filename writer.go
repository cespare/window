// Package window provides utilities for examining windows of IO streams.
package window

// A Writer is an io.Writer which retains only the last N bytes that were written.
type Writer struct {
	buf []byte // fixed ring buffer
	// i points past the last byte in the ring
	// (so if wrapped == true, then buf[i] is the first byte of the window).
	i       int
	wrapped bool
}

// NewWriter creates a Writer of the given window size.
func NewWriter(size int) *Writer {
	return &Writer{buf: make([]byte, size)}
}

// Write implements the io.Writer interface.
func (w *Writer) Write(b []byte) (int, error) {
	n := len(b)
	if n > len(w.buf) {
		b = b[n-len(w.buf):]
	}
	copy(w.buf[w.i:], b)
	rem := len(b) - len(w.buf) + w.i
	if rem > 0 {
		w.wrapped = true
		copy(w.buf, b[len(b)-rem:])
		w.i = rem
	} else {
		w.i += len(b)
	}
	return n, nil
}

// Bytes returns the current window contents.
// The returned value is a copy and is not modified by future writes.
func (w *Writer) Bytes() []byte {
	if w.wrapped {
		buf := make([]byte, len(w.buf))
		copy(buf, w.buf[w.i:])
		copy(buf[len(buf)-w.i:], w.buf[:w.i])
		return buf
	}
	buf := make([]byte, w.i)
	copy(buf, w.buf[:w.i])
	return buf
}
