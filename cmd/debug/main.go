package main

import (
	"log"
	"os"

	cb "github.com/therealfakemoot/copy-bot"
)

func main() {
	c := cb.NewChain(2)

	f, err := os.Open("corpus/kjbshort.txt")
	if err != nil {
		log.Fatalf("couldn't open corpus: %s\n", err)
	}
	defer f.Close()
	c.Build(f)
}
