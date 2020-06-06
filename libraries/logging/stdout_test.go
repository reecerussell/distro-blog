package logging

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestStdOut(t *testing.T) {
	s := StdOut()
	if s == nil {
		t.Errorf("expected an inatnce of StdOut but got nil")
	}
}

func TestStdOutWrite(t *testing.T) {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	tm := "Hello World!"
	StdOut().Write(tm)
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	str := string(buf.Bytes())

	if !strings.Contains(str, tm) {
		t.Errorf("expected os.Stdout to contain '%s'", tm)
	}
}
