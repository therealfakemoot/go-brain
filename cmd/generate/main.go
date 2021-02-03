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

	cb "github.com/therealfakemot/copy-bot"
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

func wiggle(low, high int) int {
	return rand.Intn(high) + low
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

	if stat.Size() > 0 {
		dec := json.NewDecoder(brain)
		err = dec.Decode(c)
		if err != nil {
			log.Fatalf("couldn't load brain into []byte: %#v\n", err)
		}
	}
	c = cb.NewChain(order)

	return c, nil
}

func main() {
	var walkdir, brainpath string
	var order, min, max int

	flag.IntVar(&order, "order", 1, "Ordinality of Markov chains.")
	flag.IntVar(&min, "min", 30, "Minimum output word count.")
	flag.IntVar(&max, "max", 90, "Maximum output word count.")
	flag.StringVar(&walkdir, "walkdir", ".", "path to directory containing corpus data")
	flag.StringVar(&brainpath, "brainpath", "default.brain", "path to file containing a brain")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	c, err := LoadBrain(brainpath, order)
	if err != nil {
		log.Fatalf("error loading or creating brain: %#v", err)
	}
	wf := W(c)
	filepath.Walk(walkdir, wf)

	f, err := os.OpenFile(brainpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("couldn't open brain file: %#v\n", err)
	}

	fmt.Println(c.Generate(35))
}
