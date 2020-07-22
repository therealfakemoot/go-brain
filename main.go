package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	m "github.com/therealfakemoot/gomarkov"
)

// Normalize takes an input string and returns a slice of whitespace spearated words in all lowercase with all punctuation removed.
func Normalize(s string) []string {
	// This regex is neat. \p{L} means "any letter in any language". \p{Z} means "any whitespace character in any unicode language". I'm using these so the markov engine can be 100% unicode friendly and language agnostic.
	reg, _ := regexp.Compile(`[^\p{L}\p{Z}]+`)
	words := strings.Split(strings.ToLower(reg.ReplaceAllString(s, "")), " ")
	return words

}

// LoadFile returns the contents of a file as a raw string.
func LoadFile(fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// W builds a closure that fits the WalkFunc signature so you can recursively load corpus files.
func W(c *m.Chain) filepath.WalkFunc {
	wf := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			log.Printf("Unable to walk %s.\n", path)
			return err
		}

		raw, err := LoadFile(path)
		if err != nil {
			log.Printf("Unable to load %s.\n", path)
			return err
		}
		c.Add(Normalize(raw))
		return nil
	}
	return wf
}

func wiggle(low, high int) int {
	return rand.Intn(high) + low
}

func Text(c *m.Chain, x, y int) string {

	tokens := m.NGram{m.StartToken}

	for i := 0; i < wiggle(x, y); i++ {
		lastIndex := len(tokens) - 1
		var seed m.NGram
		seed = append(seed, tokens[lastIndex:lastIndex+c.Order]...)
		/*
			for _, t := range tokens[lastIndex : lastIndex+c.Order] {
								seed = append(seed, t)
											}
		*/
		n, err := c.Generate(seed)
		if err != nil {
			fmt.Printf("%s", err)
		}
		if n == " " {
			i--
		}
		tokens = append(tokens, n)
	}
	return strings.Join(tokens[1:len(tokens)-1], " ")
}

func main() {
	var fname string
	var order, min, max int

	flag.IntVar(&order, "order", 1, "Ordinality of Markov chains.")
	flag.IntVar(&min, "min", 10, "Minimum output word count.")
	flag.IntVar(&max, "max", 30, "Maximum output word count.")
	flag.StringVar(&fname, "filename", "corpus.txt", "path to file containing corpus data")

	flag.Parse()

	raw, err := LoadFile(fname)
	if err != nil {
		log.Fatalf("%s", err)
	}
	corpus := Normalize(raw)
	c := m.NewChain(order)
	c.Add(corpus)

	fmt.Println(Text(c, min, max))
}
