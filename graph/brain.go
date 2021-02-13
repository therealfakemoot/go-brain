package graph

import (
	gonum "gonum.org/v1/gonum/graph/simple"
)

// NewBrain creates a gonum.DirectedWeightedGraph powered markov brain with a key order of n
func NewBrain(n int) Brain {
	var b Brain
	b.g = gonum.NewWeightedDirectedGraph(0.0, 0.0)

	return b
}

// Brain is a gonum.DirectedWeightedGraph powered markov brain
type Brain struct {
	g *gonum.WeightedDirectedGraph
}
