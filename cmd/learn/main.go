package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

// LoadFile returns the contents of a file as a raw string.
func LoadFile(fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// LoadBrain loads a brain from a file or creates an empty one with the given order.
func LoadBrain(fn string, order int) (*cb.Chain, error) {
	var c *cb.Chain

	brain, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("couldn't open brain file: %#v\n", err)
	}
	defer brain.Close()

	stat, err := brain.Stat()
	if err != nil {
		log.Fatalf("couldn't stat brain file: %#v\n", err)
	}

	// check if the file is empty. if not, attempt to load
	if stat.Size() > 0 {
		dec := json.NewDecoder(brain)
		err = dec.Decode(c)
		if err != nil {

			err, ok := err.(*json.InvalidUnmarshalError)

			if !ok {
				return c, fmt.Errorf("couldn't load brain into []byte: %#v\n", err)
			}

			// this is happening because sometimes loading the brain throws a weird json error
			// where it only contains a pointer to a reflect.Type
			// return c, fmt.Errorf("error unmarshaling %s.%s: %w", err.Type.PkgPath(), err.Type.Name(), err)
			return c, fmt.Errorf("error unmarshaling: %w", err)

		}
	}
	c = cb.NewChain(order)

	return c, nil
}

func main() {
	var walkdir, brainpath string
	var order int

	flag.IntVar(&order, "order", 2, "Ordinality of Markov chains.")
	flag.StringVar(&walkdir, "walkdir", ".", "path to directory containing corpus data")
	flag.StringVar(&brainpath, "brainpath", "default.brain", "path to file containing a brain")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	c := cb.NewChain(order)
	wf := cb.W(c)
	filepath.Walk(walkdir, wf)

	// f, err := os.OpenFile(brainpath, os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// log.Fatalf("couldn't open brain file: %#v\n", err)
	// }
	enc := json.NewEncoder(os.Stdout)
	err := enc.Encode(c)
	if err != nil {
		log.Fatalf("error encoding brain: %#v\n", err)
	}

	// fmt.Printf("%#v\n", c)
	// fmt.Println(c.Generate(15))
}
