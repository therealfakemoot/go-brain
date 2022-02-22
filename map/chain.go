package mapbrain

import (
	"bufio"
	// "fmt"
	"io"
	// "log"
	"math/rand"
	"regexp"
	"strings"

	"github.com/therealfakemoot/copy-bot"
)

// Normalize takes an input string and returns a slice of whitespace spearated words in all lowercase with all punctuation removed.
func Normalize(entries []string) []string {
	var res []string
	// This regex is neat. \p{L} means "any letter in any language". \p{Z} means "any whitespace character in any unicode language". I'm using these so the markov engine can be 100% unicode friendly and language agnostic.
	reg := regexp.MustCompile(`[^\p{L}\p{Z}]+`)
	for _, e := range entries {
		split := reg.Split(e, -1)
		for _, w := range split {
			res = append(res, w)
		}
		// res[0] = reg.ReplaceAllString(
		//words := strings.Split(strings.ToLower(reg.ReplaceAllString(s, "")), " ")
	}
	return res

}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	Chain     map[string][]string
	PrefixLen int
}

// NewBrain returns a new Chain with prefixes of prefixLen words.
func NewBrain(prefixLen int) Chain {
	return Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Learn(r io.Reader) {
	scanner := bufio.NewScanner(r)
	p := make(brain.Prefix, c.PrefixLen)
	for scanner.Scan() {
		for _, s := range strings.Fields(scanner.Text()) {
			key := p.String()
			c.Chain[key] = append(c.Chain[key], s)
			p.Shift(s)
		}
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c Chain) Generate(n int) []string {
	p := make(brain.Prefix, c.PrefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.Chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return words
}

// Chain is unimplemented
func (c *Chain) Save() error {
	return nil
}
