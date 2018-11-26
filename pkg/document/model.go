package document

import (
	"fmt"
	"github.com/doout/cps842/pkg/utils"
	"github.com/doout/prose"
	"hash/fnv"

	"math"
	"sort"
)

//Use for vector space model for now. Will change it to be breaking up into parts like MapReduce.
//Will use goroute as worker but in the future this can be deploy on k8s.
//Overhead will need to be taking into account when testing it

type Term struct {
	Index     uint32 // The hash of the word.
	Frequency int64  //term f
}

type DocumentSlice []Document

type Result struct {
	CosSim   float64
	Document uint64
}

//Query can be Document
type Query Document

type Document struct {
	Layout map[string][]Term
	W      *func(ts []Term, index uint, layoutItem string) float64 `json:"-"`
}

type ModelOptions struct {
	StopWords     []string
	PorterStemmer bool
}

type Model struct {
	ModelOptions
	TokenParser            []func(token string) string `json:"-"`
	DictionaryInvert       map[uint32]string           // help to get the term from the hash, mostly for debugging
	Dictionary             map[string]map[string]*Term // term -> term_index, df
	Documents              map[int]*Document           // id -> [terms -> [{index, tf}, {index, tf} .. ]]
	TotalNumberOfDocuments uint64
	Info                   map[int]map[string]string
	WD                     *func(ts []Term, index uint, layoutItem string) float64 `json:"-"`
	WQ                     *func(ts []Term, index uint, layoutItem string) float64 `json:"-"`
}

func LoadModel(folder string) *Model {
	m := MakeModel()
	tokenParser := []func(token string) string{RemovePunctuation, ToLower}
	utils.LoadJsonFromFile(&m.ModelOptions, fmt.Sprintf("%s/%s", folder, "modelOption"))
	utils.LoadJsonFromFile(&m.Dictionary, fmt.Sprintf("%s/%s", folder, "dictionary"))
	utils.LoadJsonFromFile(&m.Info, fmt.Sprintf("%s/%s", folder, "docinfo"))
	utils.LoadJsonFromFile(&m.DictionaryInvert, fmt.Sprintf("%s/%s", folder, "dictionaryInvert"))
	utils.LoadJsonFromFile(&m.Documents, fmt.Sprintf("%s/%s", folder, "posting"))

	if len(m.StopWords) > 0 {
		tokenParser = append(tokenParser, func(token string) string {
			for _, value := range m.StopWords {
				if token == value {
					return ""
				}
			}
			return token
		})
	}
	if m.PorterStemmer {
		tokenParser = append(tokenParser, PorterStemmer)
	}
	m.TokenParser = tokenParser
	m.TotalNumberOfDocuments = uint64(len(m.Documents))
	//place the func back into the terms
	for _, k := range m.Documents {
		k.W = m.WD
	}
	return m
}

// TermFrequency map[string]Item
// Since we have the term frequency let first make the dictionary

func MakeModel() *Model {
	m := Model{
		Dictionary:             make(map[string]map[string]*Term),
		DictionaryInvert:       make(map[uint32]string),
		Documents:              make(map[int]*Document),
		TotalNumberOfDocuments: uint64(0),
	}
	//For the doc
	w := func(ts []Term, index uint, layoutItem string) float64 {
		t := m.DictionaryInvert[ts[index].Index]
		df := int64(0)
		if temp, ok := m.Dictionary[layoutItem]; ok {
			if term, ok := temp[t]; ok {
				df = term.Frequency
			}
		}
		if df == 0 || ts[index].Frequency == 0 {
			return 0
		}
		idf := math.Log10(float64(m.TotalNumberOfDocuments) / float64(df))
		tf := math.Log10(float64(ts[index].Frequency)) + 1
		//We will never have Frequency 0 but just in case

		//tf := float64(ts[index].Frequency)
		//tf := math.Log10(float64(ts[index].Frequency)) + 1
		return tf * idf
	}
	m.WD = &w
	m.WQ = &w
	return &m
}

func (m *Model) BuildQuery(query map[string]string, tokenParsers ...func(token string) string) (*Query, error) {
	a := make(map[string][]prose.Token)
	for key, value := range query {
		tokens, err := getProseToken(value)
		if err != nil {
			return nil, err
		}
		a[key] = tokens.Tokens()

	}
	termF, _ := countWord(a, tokenParsers...)
	returnQuery := Query{Layout: make(map[string][]Term), W: m.WQ}

	for key, _ := range query {
		for term, f := range termF[key] {
			returnQuery.Layout[key] = append(returnQuery.Layout[key],
				Term{
					Index:     hash(term),
					Frequency: int64(f),
				})
		}
	}

	return &returnQuery, nil
}

func (m *Model) Search(input map[string]string) []Result {
	q, _ := m.BuildQuery(input, m.TokenParser...)
	q2 := Document(*q)
	r := []Result{}
	//Only search on W for now.
	for index, doc := range m.Documents {
		res := Result{0, uint64(index)}
		for key := range q2.Layout {
			cosSim := doc.CosSim(&q2, key)
			res.CosSim += cosSim
		}
		if res.CosSim > 0 {
			r = append(r, res)
		}
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].CosSim > r[j].CosSim
	})
	index := 0
	maxIndex := len(r)
	for index < maxIndex {
		if r[index].CosSim <= 0 {
			break
		}
		index++
	}

	if index == 0 {
		return []Result{}
	}

	return r[:index]
}

func (m *Model) AddDocuments(tfs map[string]map[string]Item) {
	for layoutKey, _ := range tfs {
		if _, ok := m.Dictionary[layoutKey]; !ok {
			m.Dictionary[layoutKey] = make(map[string]*Term)
		}
		for key, value := range tfs[layoutKey] {
			termIndex := hash(key)
			m.DictionaryInvert[termIndex] = key
			dictionaryTerm := &Term{Index: termIndex, Frequency: *value.DocumentFrequency}
			m.Dictionary[layoutKey][key] = dictionaryTerm
			for docId, docValue := range value.DocumentInfo {
				if modelDocValue, ok := m.Documents[docId]; ok {
					modelDocValue.Layout[layoutKey] = append(modelDocValue.Layout[layoutKey], Term{Index: termIndex, Frequency: int64(docValue.Frequency)})
				} else {
					d := Document{Layout: make(map[string][]Term), W: m.WD}
					d.Layout[layoutKey] = []Term{{Index: termIndex, Frequency: int64(docValue.Frequency)}}
					m.Documents[docId] = &d

				}
			}
		}
	}
	m.TotalNumberOfDocuments = uint64(len(m.Documents))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
