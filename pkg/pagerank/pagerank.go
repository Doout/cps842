package pagerank

import (
	"gonum.org/v1/gonum/floats"
	"math"
)

//page rank alg can be found here https://en.wikipedia.org/wiki/PageRank
type pageRank struct {
	NumberOfNodes        uint64
	InBoundLink          [][]uint64
	NumberOfOutBoundLink []int
	KeyToIndex           map[uint64]uint64
	IndexToKey           map[uint64]uint64
}

type PageRankResult struct {
	R          []float64 // page rank
	KeyToIndex map[uint64]uint64
}

func New() pageRank {
	return pageRank{
		NumberOfNodes: uint64(0),
		KeyToIndex:    make(map[uint64]uint64),
		IndexToKey:    make(map[uint64]uint64),
	}
}

func (rank *pageRank) Rank(dampingFactor, tolerance float64) PageRankResult {
	//Just do one round.
	r := make([]float64, rank.NumberOfNodes)
	index := uint64(0)
	for index < rank.NumberOfNodes {
		r[index] = float64(1) / float64(rank.NumberOfNodes)
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
	ret := PageRankResult{KeyToIndex: make(map[uint64]uint64),
		R: r}

	for key, value := range rank.KeyToIndex {
		ret.KeyToIndex[key] = value
	}
	return ret
}

func (rank *PageRankResult) Norm() {
	xMax := floats.Max(rank.R)
	xMin := floats.Min(rank.R)
	d := xMax - xMin
	r := make([]float64, len(rank.R))
	for index, value := range rank.R {
		r[index] = (value - xMin) / d
	}
	rank.R = r
}

func (rank *PageRankResult) GetPageRank(page int) float64 {
	return rank.GetPageRank64(uint64(page))
}

func (rank *PageRankResult) GetPageRank64(page uint64) float64 {
	if val, ok := rank.KeyToIndex[page]; ok {
		return rank.R[val]
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
	for index < rank.NumberOfNodes {
		links := rank.InBoundLink[index]
		sum := float64(0)
		for _, link := range links {
			sum += P[link] / float64(rank.NumberOfOutBoundLink[link])
		}

		newP := (1-dampingFactor)/float64(rank.NumberOfNodes) + (dampingFactor * (sum + linkToAll))
		P2[index] = newP
		vsum += newP
		index++
	}
	//Let make sure it add up to 1
	ivsum := 1 / vsum
	for index, val := range P2 {
		P2[index] = val * ivsum
	}
	return P2
}

//
func (pr *pageRank) getDanglingNodes() []uint64 {
	danglingNodes := make([]uint64, 0, pr.NumberOfNodes)
	for i, outBoundLinks := range pr.NumberOfOutBoundLink {
		if outBoundLinks == 0 {
			danglingNodes = append(danglingNodes, uint64(i))
		}
	}
	return danglingNodes
}

func (rank *pageRank) getIndex(docID uint64) uint64 {
	var returnIndex uint64
	var ok bool
	if returnIndex, ok = rank.KeyToIndex[docID]; !ok {
		returnIndex = rank.NumberOfNodes
		rank.KeyToIndex[docID] = returnIndex
		rank.IndexToKey[returnIndex] = docID
		rank.InBoundLink = append(rank.InBoundLink, []uint64{})
		rank.NumberOfOutBoundLink = append(rank.NumberOfOutBoundLink, 0)
		rank.NumberOfNodes++
	}
	return returnIndex
}

func (rank *pageRank) Link64(from, to uint64) {
	//add a link to the node where it in bound. A -> B mean at B add a inbound in for A
	//Check if the node exist.
	//Update from and to with the index it save
	from, to = rank.getIndex(from), rank.getIndex(to)
	for _, val := range rank.InBoundLink[to] {
		if val == from {
			//Only count one link from the doc
			return
		}
	}

	rank.InBoundLink[to] = append(rank.InBoundLink[to], from)
	rank.NumberOfOutBoundLink[from]++

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
