package logging

import (
	"io"
	"os"
)

// stdOut is a logger output implementation used to write to the stdout.
type stdOut struct{}

// StdOut returns a new Output.
func StdOut() Output {
	return &stdOut{}
}

// Write writes a message to the stdout.
func (*stdOut) Write(message string) {
	io.WriteString(os.Stdout, message)
}
