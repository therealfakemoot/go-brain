package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"

	cb "github.com/therealfakemoot/copy-bot"
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

func main() {
	var corpus, brainpath string
	var order int

	flag.IntVar(&order, "order", 2, "Ordinality of Markov chains.")
	flag.StringVar(&corpus, "corpus", ".", "path to directory containing corpus data")
	flag.StringVar(&brainpath, "brainpath", "default.brain", "path to file containing a brain")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	c := cb.NewChain(order)
	wf := cb.W(c)
	filepath.Walk(corpus, wf)

	f, err := os.OpenFile(brainpath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("couldn't open brain file: %#v\n", err)
	}
	enc := json.NewEncoder(f)
	err = enc.Encode(c)
	if err != nil {
		log.Fatalf("error encoding brain: %#v\n", err)
	}
}
