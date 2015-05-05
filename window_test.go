package window

import "testing"

func TestWriter(t *testing.T) {
	w := NewWriter(4)
	w.Write([]byte("a"))
	if got, want := string(w.Bytes()), "a"; got != want {
		t.Fatalf("got %q; want: %q", got, want)
	}
	w.Write([]byte("bc"))
	if got, want := string(w.Bytes()), "abc"; got != want {
		t.Fatalf("got %q; want: %q", got, want)
	}
	w.Write([]byte("d"))
	if got, want := string(w.Bytes()), "abcd"; got != want {
		t.Fatalf("got %q; want: %q", got, want)
	}
	w.Write([]byte("e"))
	if got, want := string(w.Bytes()), "bcde"; got != want {
		t.Fatalf("got %q; want: %q", got, want)
	}
	w.Write([]byte("f"))
	if got, want := string(w.Bytes()), "cdef"; got != want {
		t.Fatalf("got %q; want: %q", got, want)
	}
	w.Write([]byte("foobarfoobar"))
	if got, want := string(w.Bytes()), "obar"; got != want {
		t.Fatalf("got %q; want: %q", got, want)
	}
}
