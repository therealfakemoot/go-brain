package chain

import (
	"fmt"
	"math/rand"
	"strings"

	m "github.com/therealfakemoot/gomarkov"
)

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
