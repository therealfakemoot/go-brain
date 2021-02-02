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
	"strings"

	m "github.com/therealfakemoot/gomarkov"
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

		// raw, err := LoadFile(path)
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			log.Printf("Unable to open file %s: %s\n", path, err)
			return err
		}

		var (
			s       string
			entries []string
		)

		for {
			if _, err := fmt.Fscan(f, &s); err != nil {
				break
			}
			entries = append(entries, s)
		}
		c.Add(Normalize(entries))
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
	var walkdir, brainpath string
	var order, min, max int

	flag.IntVar(&order, "order", 1, "Ordinality of Markov chains.")
	flag.IntVar(&min, "min", 30, "Minimum output word count.")
	flag.IntVar(&max, "max", 90, "Maximum output word count.")
	flag.StringVar(&walkdir, "walkdir", ".", "path to directory containing corpus data")
	flag.StringVar(&brainpath, "brainpath", "", "path to file containing a brain")

	flag.Parse()

	var c *m.Chain

	if brainpath != "" {
		brain, err := os.Open(brainpath)
		defer brain.Close()
		if err == os.ErrNotExist {
			log.Fatalf("brain doesn't exist for loading: %s", err)
		}
		if err != nil {
			log.Fatalf("error loading brain file: %s", err)
		}

		dec := json.NewDecoder(brain)
		err = dec.Decode(c)
		if err != nil {
			log.Fatalf("couldn't load brain into []byte: %s", err)
		}
	}
	c = m.NewChain(order)

	wf := W(c)
	filepath.Walk(walkdir, wf)
	/*
		raw, err := LoadFile(fname)
		if err != nil {
			log.Fatalf("%s", err)
		}
		corpus := Normalize(raw)
		c.Add(corpus)
	*/

	fmt.Println(Text(c, min, max))
	f, err := os.OpenFile(brainpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("couldn't open brain file for dumping: %s", err)
	}
	enc := json.NewEncoder(f)
	enc.Encode(c)
}
