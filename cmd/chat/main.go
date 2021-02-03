package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	cb "github.com/therealfakemoot/copy-bot"
)

func main() {
	var (
		brain  string
		length int
		b      cb.Chain
	)

	flag.IntVar(&length, "length", 25, "desired message length in words")
	flag.StringVar(&brain, "brain", "default.brain", "path to brain file")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	brainFile, err := os.Open(brain)
	if err != nil {
		log.Fatalf("could not open brain jar: %s", err)
	}
	dec := json.NewDecoder(brainFile)
	err = dec.Decode(&b)
	if err != nil {
		log.Fatalf("could not retrieve brain from jar: %s", err)
	}

	fmt.Println(b.Generate(length))
}
