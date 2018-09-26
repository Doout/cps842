package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"sync/atomic"
)

type Item struct {
	TotalFrequency *int64
	DocumentInfo   frequencyMap
}

type DocumentInfo struct {
	Frequency int
	Location  []int
}

type frequencyMap map[int]*DocumentInfo

//Recreate the JSON func use for json.Marshal since map was not getting sort by number.
func (t frequencyMap) MarshalJSON() (text []byte, err error) {
	e := bytes.Buffer{}
	e.WriteByte('{')
	// Extract and sort the keys.
	sv := make([]int, len(t))
	index := 0
	for i := range t {
		sv[index] = i
		index += 1
	}
	sort.Ints(sv)
	for i, kv := range sv {
		if i > 0 {
			e.WriteByte(',')
		}
		e.WriteString(fmt.Sprintf(`"%d"`, kv))
		e.WriteByte(':')
		b, err := json.Marshal(t[kv])
		if err != nil {
			return nil, err
		}
		e.Write(b)
	}
	e.WriteByte('}')
	return e.Bytes(), nil
}

func NewItem() Item {
	total := int64(0)
	return Item{&total, make(map[int]*DocumentInfo)}
}

func (item *Item) GetTotalFrequency() int64 {
	return *item.TotalFrequency
}

func (item *Item) GetFrequency(documentId int) int {
	return item.DocumentInfo[documentId].Frequency
}

func (item *Item) AddFrequency(documentId, frequency int, location []int) {
	item.DocumentInfo[documentId] = &DocumentInfo{Frequency: frequency}
	item.DocumentInfo[documentId].Frequency = frequency
	item.DocumentInfo[documentId].Location = location
	atomic.AddInt64(item.TotalFrequency, int64(frequency))
}
