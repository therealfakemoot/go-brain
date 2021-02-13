package graph

import (
	"bufio"
	"io"
	"math/rand"
	"strings"

	"github.com/therealfakemoot/copy-bot"

	ggraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// Brain is a gonum.DirectedWeightedGraph powered markov brain
type Brain struct {
	o int // how many consecutive words to use as each vertex
	g *simple.WeightedDirectedGraph
}

// NewBrain creates a gonum.DirectedWeightedGraph powered markov brain with a key order of n
func NewBrain(n int) Brain {
	var b Brain
	b.o = n
	b.g = simple.NewWeightedDirectedGraph(0.0, 0.0)

	return b
}

// Generate produces up to n strings from the stored brain juice
func (b *Brain) Generate(n int) []string {

	var s []string

	weighted := make([]string, 1000)

	// create the baseline start token, n empty strings
	p := make(brain.Prefix, b.o)
	n1 := Node{p.String()}

	var matches []ggraph.WeightedEdge
	var loopNode Node = n1
	for len(s) < n {
		matches = EdgesOf(b.g, loopNode)
		for _, m := range matches {
			w := m.Weight()
			to := m.To().(Node)
			for i := 0.0; i < w; i++ {
				weighted = append(weighted, to.Key)
			}
		}
		choice := weighted[rand.Intn(len(weighted))]
		s = append(s, choice)
		p.Shift(choice)
		loopNode = Node{p.String()}
		weighted = weighted[:]

	}

	return s
}

func (b *Brain) Learn(r io.Reader) {
	scanner := bufio.NewScanner(r)
	p := make(brain.Prefix, b.o)
	for scanner.Scan() {
		for _, s := range strings.Fields(scanner.Text()) {
			n1 := Node{p.String()}
			p.Shift(s)
			n2 := Node{p.String()}
			weight, ok := b.g.Weight(n1.ID(), n2.ID())
			if !ok {
				we := simple.WeightedEdge{
					F: n1,
					T: n2,
					W: 1.0,
				}
				b.g.SetWeightedEdge(we)
				continue
			}

			we := simple.WeightedEdge{
				F: n1,
				T: n2,
				W: weight + 1.0,
			}
			b.g.SetWeightedEdge(we)
		}
	}
}

func (b *Brain) Save() error {

	return nil
}
