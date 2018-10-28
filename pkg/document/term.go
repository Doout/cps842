package document

import (
	"math"
	"sort"
)

//this file hold all the func for Terms
func (t *Document) Dot(t2 *Document, layoutItem string) float64 {
	terms1 := t.Layout[layoutItem]
	terms2 := t2.Layout[layoutItem]
	sort.Slice(terms1, func(i, j int) bool {
		return terms1[i].Index < terms1[j].Index
	})
	sort.Slice(terms2, func(i, j int) bool {
		return terms2[i].Index < terms2[j].Index
	})
	i, j := 0, 0
	iMax, jMax := len(terms1), len(terms2)
	total := float64(0)
	//Check if the term match, if yes time tf and move to the next item in the list.
	//If i or j go over the max we can assume everything else will be zero
	for i < iMax && j < jMax {
		//check if the same
		if terms1[i].Index == terms2[j].Index {
			w := *t.W
			w2 := *t2.W
			total += float64(w(terms1, uint(i), layoutItem) * w2(terms2, uint(j), layoutItem))
			i++
			j++
		} else if terms1[i].Index > terms2[j].Index {
			//Since this is sorts we need to get a higher number of t2
			j++
		} else if terms1[i].Index < terms2[j].Index {
			//Since t2 is bigger we need to check if the next item in t match
			i++
		}
	}
	return total
}

//
func (t *Document) CosSim(t2 *Document, LayoutItem string) float64 {
	a := t.Dot(t2, LayoutItem) / (t.Len(LayoutItem) * t2.Len(LayoutItem))
	if math.IsNaN(a) {
		return 0
	} else {
		return a
	}
}

func (t *Document) Len(layoutItem string) float64 {
	l := float64(0)
	for _, term := range (*t).Layout[layoutItem] {
		l += float64(term.Frequency * term.Frequency)
	}
	return math.Sqrt(l)
}
