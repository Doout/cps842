package main

import (
	_ "github.com/dcadenas/pagerank"
	g "github.com/doout/cps842/pkg/pagerank"
)

func main() {
	graph := g.New()
	//graph2 := g.New()

	A := 0
	B := 1
	C := 2
	D := 3
	E := 4
	F := 5

	graph.Link(B, C)
	graph.Link(C, B)
	graph.Link(D, A)
	graph.Link(D, B)
	graph.Link(E, B)
	graph.Link(E, D)
	graph.Link(E, F)
	graph.Link(F, B)
	graph.Link(F, E)
	graph.Link(10, B)
	graph.Link(10, E)
	graph.Link(11, B)
	graph.Link(11, E)
	graph.Link(12, B)
	graph.Link(12, E)
	graph.Link(13, E)
	graph.Link(14, E)

	graph.Link(1, 2)

	//17

	a := graph.Rank(0.85, 0.00001)
	//Set range from 0 to 1.
	a.Norm()

}
