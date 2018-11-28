package document

import (
	"github.com/doout/prose"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var vaildToken = regexp.MustCompile(`(?m)^(\d+|\w+)$`)

var DocumentsItems = []string{"T", "W", "N", "A"}
var DocumentLink = "X"

type TermFrequencys struct {
	TermFrequency map[string]map[string]Item
	Info          map[int]map[string]string
}

func (doc *TermFrequencys) GetDocumentsData() map[int]map[string]string {
	return doc.Info
}

func BuildDocument(b []map[string]string) *TermFrequencys {
	return BuildDocumentWithTokenParser(b)
}

func BuildDocumentWithTokenParser(b []map[string]string, tokenParsers ...func(token string) string) *TermFrequencys {
	d := TermFrequencys{
		TermFrequency: make(map[string]map[string]Item),
		Info:          make(map[int]map[string]string),
	}

	type tokensChanData struct {
		ID int
		// itemName -> word -> count
		Tokens      map[string]map[string]int
		Info        map[string]string
		occurrences map[string]map[string][]int
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
			if _, ok := d.TermFrequency[key]; !ok {
				d.TermFrequency[key] = make(map[string]Item)
			}
			for word, token := range value {
				if v, ok := d.TermFrequency[key][word]; ok {
					v.AddFrequency(wordMap.ID, token)
				} else {
					v = NewItem()
					v.AddFrequency(wordMap.ID, token)
					d.TermFrequency[key][word] = v
				}
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
func GetToken(s map[string]string, items ...string) map[string][]prose.Token {
	returnTokens := make(map[string][]prose.Token, len(items))

	type Tokens struct {
		Tokens   []prose.Token
		ItemName string
	}
	tokensThreadChan := make(chan Tokens, len(items))
	for _, item := range items {
		tempItems := item
		//tempIndex := i
		//Per each strings get the tokens
		go func() {
			doc, err := getProseToken(s[tempItems])
			if err != nil {
				log.Fatal(err)
			}
			//baseOffset := 32 - uint(math.Log2(float64(len(items))))
			//tokens2 := make([]prose.Token, len(doc.Tokens()))
			////Update the index to match the location in which this token can be found in
			//for index, v := range doc.Tokens() {
			//	v.Index |= tempIndex << baseOffset
			//	tokens2[index] = v
			//}
			//Send the tokens to the main thread to be sync
			tokensThreadChan <- Tokens{doc.Tokens(), tempItems}
		}()
	}

	//Sync the background thread this function spin up
	for _ = range items {
		tok := <-tokensThreadChan
		returnTokens[tok.ItemName] = tok.Tokens
	}
	return returnTokens
}

/*
Return the lists of words with the index of the word in the documnt
This function does not find the index of the word itself but use token.Index
*/
func countWord(tokens map[string][]prose.Token, tokenParsers ...func(token string) string) (map[string]map[string]int, map[string]map[string][]int) {
	//itemName -> word -> count
	li := make(map[string]map[string]int)
	//itemName -> word -> location
	oc := make(map[string]map[string][]int)

	for key, value := range tokens {
		li[key] = make(map[string]int)
		oc[key] = make(map[string][]int)
	OUTER:
		for _, token := range value {
			word := token.Text
			for _, fn := range tokenParsers {
				if strings.Compare(word, "") == 0 {
					continue OUTER
				}
				word = fn(word)
			}
			if _, ok := li[key][word]; ok {
				li[key][word] += 1
				oc[key][word] = append(oc[key][word], token.Index)
			} else {
				li[key][word] = 1
				oc[key][word] = []int{token.Index}
			}
		}
	}
	return li, oc
}
