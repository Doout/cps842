package pagerank

import (
	"gonum.org/v1/gonum/floats"
	"math"
)

//page rank alg can be found here https://en.wikipedia.org/wiki/PageRank
type pageRank struct {
	numberOfLinks        uint64
	inBoundLink          [][]uint64
	numberOfOutBoundLink []int
	keyToIndex           map[uint64]uint64
	indexToKey           map[uint64]uint64
}
type pageRankResult struct {
	r          []float64 // page rank
	keyToIndex map[uint64]uint64
}

func New() pageRank {
	return pageRank{
		numberOfLinks: uint64(0),
		keyToIndex:    make(map[uint64]uint64),
		indexToKey:    make(map[uint64]uint64),
	}
}

func (rank *pageRank) Rank(dampingFactor, tolerance float64) pageRankResult {
	//Just do one round.
	r := make([]float64, rank.numberOfLinks)
	index := uint64(0)
	for index < rank.numberOfLinks {
		r[index] = float64(1) / float64(rank.numberOfLinks)
		index++
	}
	danglingNodes := rank.getDanglingNodes()

	c := tolerance + 1
	for c > tolerance {
		r2 := rank.step(dampingFactor, r, danglingNodes)
		c = calculateChange(r, r2)
		r = r2
	}

	// Create the target map
	ret := pageRankResult{keyToIndex: make(map[uint64]uint64),
		r: r}
	for key, value := range rank.keyToIndex {
		ret.keyToIndex[key] = value
	}

	return ret
}

func (rank *pageRankResult) Norm() {
	xMax := floats.Max(rank.r)
	xMin := floats.Min(rank.r)
	d := xMax - xMin
	r := make([]float64, len(rank.r))
	for index, value := range rank.r {
		r[index] = (value - xMin) / d
	}
	rank.r = r
}

func (rank *pageRankResult) GetPageRank(page int) float64 {
	return rank.GetPageRank64(uint64(page))
}

func (rank *pageRankResult) GetPageRank64(page uint64) float64 {
	if val, ok := rank.keyToIndex[page]; ok {
		return rank.r[val]
	} else {
		return 0.0
	}

}

func (rank *pageRank) step(dampingFactor float64, P []float64, danglingNodes []uint64) []float64 {
	vsum := float64(0)
	index := uint64(0)

	danglingNodesSum := float64(0)
	for _, PIndex := range danglingNodes {
		danglingNodesSum += P[PIndex]
	}

	linkToAll := danglingNodesSum / float64(len(P))

	P2 := make([]float64, len(P))
	for index < rank.numberOfLinks {
		links := rank.inBoundLink[index]
		sum := float64(0)
		for _, link := range links {
			sum += P[link] / float64(rank.numberOfOutBoundLink[link])
		}
		newP := (1-dampingFactor)/float64(rank.numberOfLinks) + (dampingFactor * (sum + linkToAll))
		P2[index] = newP
		vsum += newP
		index++
	}

	return P2
}

//
func (pr *pageRank) getDanglingNodes() []uint64 {
	danglingNodes := make([]uint64, 0, pr.numberOfLinks)

	for i, outBoundLinks := range pr.numberOfOutBoundLink {
		if outBoundLinks == 0 {
			danglingNodes = append(danglingNodes, uint64(i))
		}
	}

	return danglingNodes
}

func (rank *pageRank) getIndex(index uint64) uint64 {
	var returnIndex uint64
	var ok bool
	if returnIndex, ok = rank.keyToIndex[index]; !ok {
		returnIndex = rank.numberOfLinks
		rank.keyToIndex[index] = returnIndex
		rank.indexToKey[returnIndex] = index
		rank.inBoundLink = append(rank.inBoundLink, []uint64{})
		rank.numberOfOutBoundLink = append(rank.numberOfOutBoundLink, 0)
		rank.numberOfLinks++
	}
	return returnIndex
}

func (rank *pageRank) Link64(from, to uint64) {
	//add a link to the node where it in bound. A -> B mean at B add a inbound in for A
	//Check if the node exist.
	//Update from and to with the index it save
	from, to = rank.getIndex(from), rank.getIndex(to)
	rank.inBoundLink[to] = append(rank.inBoundLink[to], from)
	rank.numberOfOutBoundLink[from]++
}

func (rank *pageRank) Link(from, to int) {
	rank.Link64(uint64(from), uint64(to))
}

func calculateChange(p, new_p []float64) float64 {
	acc := 0.0
	for i, pForI := range p {
		acc += math.Abs(pForI - new_p[i])
	}
	return acc
}
