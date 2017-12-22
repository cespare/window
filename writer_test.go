package window

import "testing"

func TestWriter(t *testing.T) {
	w := NewWriter(4)
	for _, tt := range []struct {
		s    string
		want string
	}{
		{"a", "a"},
		{"bc", "abc"},
		{"d", "abcd"},
		{"e", "bcde"},
		{"f", "cdef"},
		{"foobarfoobar", "obar"},
	} {
		if _, err := w.Write([]byte(tt.s)); err != nil {
			t.Fatal(err)
		}
		if got := string(w.Bytes()); got != tt.want {
			t.Fatalf("after writing %q, got %q; want %q", tt.s, got, tt.want)
		}
	}
}
