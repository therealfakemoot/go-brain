package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	gb "github.com/therealfakemoot/copy-bot/graph"
)

func main() {
	var (
		brain  string
		length int
		b      gb.Brain
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
