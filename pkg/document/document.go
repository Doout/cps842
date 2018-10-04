package document

import (
	"bytes"
	"fmt"
	"github.com/bclicn/color"
	"github.com/doout/prose"
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var vaildToken = regexp.MustCompile(`(?m)^(\d+|\w+)$`)

var DocumentsItems = []string{"T", "W"}

type Documents struct {
	TermFrequency map[string]Item
	Info          map[int]map[string]string
}

func (doc *Documents) GetDocumentsData() map[int]map[string]string {
	return doc.Info
}

func (doc *Documents) GetPosting() map[string]Item {
	return doc.TermFrequency
}

func (doc *Documents) GetDictionarySort() []string {
	li := make([]string, len(doc.TermFrequency))
	index := 0
	for key, value := range doc.TermFrequency {
		li[index] = fmt.Sprintf("%s, df:%d", key, len(value.DocumentInfo))
		index++
	}
	sort.Strings(li)
	return li
}

func BuildDocument(b []map[string]string) *Documents {
	return BuildDocumentWithTokenParser(b)
}
func BuildDocumentWithTokenParser(b []map[string]string, tokenParsers ...func(token string) string) *Documents {
	d := Documents{
		TermFrequency: make(map[string]Item),
		Info:          make(map[int]map[string]string),
	}

	type tokensChanData struct {
		ID          int
		Tokens      map[string]int
		Info        map[string]string
		occurrences map[string][]int
	}
	tokensChan := make(chan tokensChanData)
	//Create a goroutine per each doc and grab all the tokens.
	for _, item := range b {
		itemTemp := item
		go func() {
			id, err := strconv.Atoi(itemTemp["I"])
			if err != nil {
				panic(err)
			}
			t := GetToken(itemTemp, DocumentsItems...)
			wordMap, wordIndexs := countWord(t, tokenParsers...)
			info := make(map[string]string, len(DocumentsItems))
			for _, value := range DocumentsItems {
				if v, ok := itemTemp[value]; ok {
					info[value] = v
				}
			}
			tcd := tokensChanData{id, wordMap, info, wordIndexs}
			//Send the data to tokensChan for it to be sync.
			tokensChan <- tcd
		}()
	}
	index := 0
	max := len(b)
	for index < max {
		//Get the data from the chan, This will wait for data to be push in
		wordMap := <-tokensChan
		d.Info[wordMap.ID] = wordMap.Info
		//For every word add them to the doc struct
		for key, value := range wordMap.Tokens {
			if v, ok := d.TermFrequency[key]; ok {
				v.AddFrequency(wordMap.ID, value, wordMap.occurrences[key])
			} else {
				v := NewItem()
				v.AddFrequency(wordMap.ID, value, wordMap.occurrences[key])
				d.TermFrequency[key] = v
			}
		}
		index++
	}
	return &d
}

//Get the tokens from the string
func getProseToken(data string) (*prose.Document, error) {
	return prose.NewDocument(data,
		prose.WithExtraction(false),
		prose.WithTagging(false),
		prose.WithSegmentation(false),
		prose.WithTokenization(true))
}

//Return a list of token from a string and update the index to match the index.
func GetToken(s map[string]string, items ...string) []prose.Token {
	var tokens []prose.Token
	tokensThreadChan := make(chan []prose.Token, len(items))
	for i, item := range items {
		tempItems := item
		tempIndex := i
		//Per each strings get the tokens
		go func() {
			doc, err := getProseToken(s[tempItems])
			if err != nil {
				log.Fatal(err)
			}
			baseOffset := 32 - uint(math.Log2(float64(len(items))))
			tokens2 := make([]prose.Token, len(doc.Tokens()))
			//Update the index to match the location in which this token can be found in
			for index, v := range doc.Tokens() {
				v.Index |= tempIndex << baseOffset
				tokens2[index] = v
			}
			//Send the tokens to the main thread to be sync
			tokensThreadChan <- tokens2
		}()
	}

	//Sync the background thread this function spin up
	for _ = range items {
		tok := <-tokensThreadChan
		tokens = append(tokens, tok...)
	}
	return tokens
}

/*
Return the lists of words with the index of the word in the documnt
This function does not find the index of the word itself but use token.Index
*/
func countWord(tokens []prose.Token, tokenParsers ...func(token string) string) (map[string]int, map[string][]int) {
	li := make(map[string]int)
	oc := make(map[string][]int)
OUTER:
	for _, token := range tokens {
		word := token.Text
		for _, fn := range tokenParsers {
			if strings.Compare(word, "") == 0 {
				continue OUTER
			}
			word = fn(word)
		}
		if _, ok := li[word]; ok {
			li[word] += 1
			oc[word] = append(oc[word], token.Index)
		} else {
			li[word] = 1
			oc[word] = []int{token.Index}
		}
	}
	return li, oc
}

//
func (d *Documents) GetTermSum(word string) string {
	word = strings.TrimSpace(word)
	fmt.Println("Looking for ", word)
	if item, ok := d.TermFrequency[word]; ok {
		di := item.DocumentInfo
		if len(di) <= 0 {
			return ""
		}
		sv := make([]int, len(di))
		index := 0
		for i := range di {
			sv[index] = i
			index += 1
		}
		sort.Ints(sv)
		buffer := bytes.Buffer{}
		for _, docId := range sv {
			di := item.DocumentInfo[docId]
			l := di.Location[0]
			index, find := DecodeLocation(l, DocumentsItems...)
			sum := getNextXToken(d.Info[docId][find], index, 10)
			buffer.WriteString(fmt.Sprintf("Doc ID: %d\n", docId))
			buffer.WriteString(fmt.Sprintf("Doc Title: %s\n", d.Info[docId]["T"]))
			buffer.WriteString(fmt.Sprintf("Term frequency: %d\n", di.Frequency))
			buffer.WriteString(sum)
			buffer.WriteString("\n\n")
		}
		return buffer.String()
	}
	return ""
}

func getNextXToken(sent string, startIndex int, numberOfNextItem int) string {
	output := bytes.Buffer{}
	if len(sent) < startIndex {
		return ""
	}
	scan := sent[startIndex:]
	scan = strings.Replace(scan, "\n", " ", -1)
	tokens := strings.Split(scan, " ")
	col := tokens[0]
	for index, tok := range tokens {
		if index > numberOfNextItem {
			break
		}
		if col == tok {
			output.WriteString(color.Blue(tok))
		} else {
			output.WriteString(tok)
		}
		output.WriteString(" ")
	}
	return strings.TrimSpace(output.String())
}

func DecodeLocation(location int, items ...string) (int, string) {
	baseOffset := 32 - uint(math.Log2(float64(len(items))))
	base := location >> baseOffset
	index := location & ((1 << baseOffset) - 1)
	return index, items[base]
}
