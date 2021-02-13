package brain

import (
	"io"
)

// Brain captures the essential functionality of a chat brain: learning, speaking, and saving
//
// By virtue of being an interface, implementations can have further methods for manipulating state
// in between calls to Generate.
type Brain interface {
	// Learn consumes text from r and adds new entries to the markov chain.
	Learn(r io.Reader)
	// Generate uses existing markov chain data to generate up to n strings
	Generate(n int) []string
	// Save is responsible for serializing brain state
	// if an implementation uses a database or something incompatible with io.Writer,
	// you can pass in io/ioutil.Discard
	Save(w io.Writer) error
}
