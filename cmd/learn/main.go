package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	cb "github.com/therealfakemoot/go-brain"
	gb "github.com/therealfakemoot/go-brain/graph"
)

func main() {
	var corpus, brainpath string
	var order int

	flag.IntVar(&order, "order", 2, "Ordinality of Markov chains.")
	flag.StringVar(&corpus, "corpus", ".", "path to directory containing corpus data")
	flag.StringVar(&brainpath, "brainpath", "default.brain", "path to file containing a brain")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	c := gb.NewBrain(order)
	wf := cb.W(&c)
	filepath.Walk(corpus, wf)

	words := c.Generate(35)
	log.Printf("words returned: %d", len(words))
	fmt.Println(words)

	return
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
