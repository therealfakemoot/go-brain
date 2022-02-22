package graph

import (
	"fmt"
	"hash/fnv"

	ggraph "gonum.org/v1/gonum/graph"
	// "gonum.org/v1/gonum/graph/simple"
)

type Node struct {
	Key string
}

func (n Node) ID() int64 {
	h := fnv.New32()
	fmt.Fprint(h, n.Key)

	return int64(h.Sum32())
}

func (n Node) String() string {
	return n.Key
}

// EdgesOf pulls all edges which start at the given node
func EdgesOf(g *WeightedDirectedGraph, n ggraph.Node) []ggraph.WeightedEdge {
	var matches []ggraph.WeightedEdge
	edges := g.WeightedEdges()
	for edges.Next() {
		we := edges.WeightedEdge()
		if we.From().ID() == n.ID() {
			matches = append(matches, we)
		}
	}

	return matches
}
