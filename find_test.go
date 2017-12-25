package window

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func init() { minBuf = 4 }

func TestFind(t *testing.T) {
	for _, tt := range []struct {
		text string
		pat  string
	}{
		{"foo", ""},
		{"foo", "f"},
		{"foo", "o"},
		{"foobar", "foo"},
		{"foobar", "bar"},
		{"the quick brown fox jumped over the lazy dog", "the"},
		{"the quick brown fox jumped over the lazy dog", "fox"},
		{"the quick brown fox jumped over the lazy dog", "dog"},
		{"the quick brown fox jumped over the lazy dog", "zzz"},
		{"", "hello"},
		{"hello world", "hello world"},
		{"dog", "the quick brown fox jumped over the lazy dog"},
	} {
		want := int64(strings.Index(tt.text, tt.pat))
		name := fmt.Sprintf("Find(%q, %q)", tt.text, tt.pat)
		check := func(variant string, r io.Reader) bool {
			t.Helper()
			tag := name
			if variant != "" {
				tag += " (" + variant + ")"
			}
			got, err := Find(r, []byte(tt.pat))
			if err != nil {
				t.Errorf("%s: got err=%s", tag, err)
				return false
			}
			if got != want {
				t.Errorf("%s: got n=%d; want %d", tag, got, want)
				return false
			}
			return true
		}
		if !check("", strings.NewReader(tt.text)) {
			continue
		}
		check("HalfReader", iotest.HalfReader(strings.NewReader(tt.text)))
		check("OneByteReader", iotest.OneByteReader(strings.NewReader(tt.text)))
		check("DataErrReader", iotest.DataErrReader(strings.NewReader(tt.text)))
		if want >= 0 {
			check("errReader", errReader{strings.NewReader(tt.text)})
		} else {
			got, err := Find(errReader{strings.NewReader(tt.text)}, []byte(tt.pat))
			if got != -1 || err != errBoom {
				t.Errorf("%s (errReader): got (%d, %s); want (-1, errBoom)",
					name, got, err)
			}
		}
	}
}

type errReader struct {
	r io.Reader
}

var errBoom = errors.New("boom")

func (er errReader) Read(b []byte) (int, error) {
	n, err := er.r.Read(b)
	if err == io.EOF {
		err = errBoom
	}
	return n, err
}
