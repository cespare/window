package window

import (
	"bytes"
	"io"
)

var minBuf = 4096 // altered in tests

// Find locates the index of the first occurrence of pat in r.
// Even if it successfully locates the search sequence,
// it may read arbitrarily far in the reader.
// If Find locates pat, it returns i >= 0, nil.
// If Find does not locate pat before reading io.EOF, it returns -1, nil.
// If Find encounters some non-EOF error err before locating pat, it returns -1, err.
func Find(r io.Reader, pat []byte) (int64, error) {
	m := len(pat)
	if m == 0 {
		return 0, nil
	}

	// Use the Boyer-Moore-Horspool algorithm for simplicity.
	var skip [256]int
	for i := range skip {
		skip[i] = m
	}
	for i, c := range pat[:m-1] {
		skip[c] = m - i - 1
	}

	bufSize := minBuf
	if m*2 > bufSize {
		bufSize = m * 2
	}
	var (
		b     = make([]byte, 0, bufSize)
		i     int   // into b
		iprev int64 // previously read count
		n     int   // Read result
		err   error
	)

	for {
		for len(b)-i < m {
			if err != nil {
				if err == io.EOF {
					return -1, nil
				}
				return -1, err
			}
			if cap(b)-i < m {
				// Slide down.
				copy(b, b[i:])
				b = b[:len(b)-i]
				iprev += int64(i)
				i = 0
			}
			n, err = r.Read(b[len(b):cap(b)])
			b = b[:len(b)+n]
		}
		c := b[i+m-1]
		if pat[m-1] == c && bytes.Equal(b[i:i+m-1], pat[:m-1]) {
			return iprev + int64(i), nil
		}
		i += skip[c]
	}
}
