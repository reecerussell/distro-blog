package normalization

import "testing"

var testStrings = map[string]string{
	"hello":           "HELLO",
	"world":           "WORLD",
	"hello world!":    "HELLO WORLD!",
	"my name is john": "MY NAME IS JOHN",
}

func TestNoramlise(t *testing.T) {
	nmzr := New()

	for in, expected := range testStrings {
		if nmzr.Normalize(in) != expected {
			t.Errorf("expected '%s' to normalize to '%s'", in, expected)
		}
	}
}
