package normalization

import "strings"

// Normalizer is a service interface use to normalize data.
type Normalizer interface {
	Normalize(in string) string
}

// New returns a new implementation of Normalizer.
func New() Normalizer {
	return new(normalizer)
}

// Normalizer implementation.
type normalizer struct{}

// Normalize normalises a given string input.
func (*normalizer) Normalize(in string) string {
	return strings.ToUpper(in)
}
