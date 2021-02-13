package brain

import (
	"io"
)

// Brain captures the essential functionality of a chat brain: learning and speaking
//
type Brain interface {
	// Learn consumes text from r and adds new entries to the markov chain.
	Learn(r io.Reader)
	// Generate uses existing markov chain data to generate up to n strings
	Generate(n int) []string
	Save() error
}
