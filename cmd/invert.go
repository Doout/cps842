package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/doout/cps842/pkg/document"
	"github.com/doout/cps842/pkg/pagerank"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var tokenParser []func(token string) string
var (
	inputFile    = ""
	outputFolder = ""
	re           = regexp.MustCompile(`(?m)^\.[A-Z]($| \d*$)`)
	stopLimit    = ""
	stopword     []string
	lower        = false
	porter       = false
)

// inv represents the playbook command
var inv = &cobra.Command{
	Use:   "invert",
	Short: "Generate inverted index from CACM collection",
	Long:  `Take a collection of documents and generate its inverted index.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		data := loadFile(inputFile)
		var doc *document.TermFrequencys
		//Remove punctuation that we don't want.
		tokenParser := []func(token string) string{document.RemovePunctuation, document.ToLower}
		//Add the stop limit token parse
		if stopLimit != "" {
			dat, err := ioutil.ReadFile(stopLimit)
			if err != nil {
				panic(err)
			}
			stopword = strings.Split(string(dat), "\n")
			tokenParser = append(tokenParser, func(token string) string {
				for _, value := range stopword {
					if token == value {
						return ""
					}
				}
				return token
			})
		}
		//Add the Porter Stemmer token parse
		if porter {
			tokenParser = append(tokenParser, document.PorterStemmer)
		}
		if tokenParser != nil {
			doc = document.BuildDocumentWithTokenParser(data, tokenParser...)
		} else {
			doc = document.BuildDocument(data)
		}
		end := time.Now()
		fmt.Println(end.Sub(start), "to process file")
		_ = doc

		type doclink struct {
			from uint64
			to   uint64
		}

		doclinks := []doclink{}

		for _, item := range data {
			links_string := item["X"]
			links := strings.Split(links_string, "\n")
			for _, link := range links {
				data := strings.Split(link, "\t")
				if data[1] == "5" {
					// X -> Y, X have a link to Y
					from, _ := strconv.Atoi(data[0])
					to, _ := strconv.Atoi(data[2])
					doclinks = append(doclinks, doclink{from: uint64(from), to: uint64(to)})
				}
			}
		}

		graph := pagerank.New()

		for _, item := range doclinks {
			graph.Link64(item.from, item.to)
		}

		scores := graph.Rank(0.85, 1e-100)

		model := document.MakeModel()
		model.StopWords = stopword
		model.TokenParser = tokenParser
		model.PorterStemmer = porter
		model.AddDocuments(doc.TermFrequency)
		model.Info = doc.Info

		start = time.Now()
		saveFile(model.Dictionary, fmt.Sprintf("%s/%s", outputFolder, "dictionary"))
		saveFile(graph, fmt.Sprintf("%s/%s", outputFolder, "pagerank_graph"))
		saveFile(scores, fmt.Sprintf("%s/%s", outputFolder, "pagerank_scores"))
		saveFile(model.Info, fmt.Sprintf("%s/%s", outputFolder, "docinfo"))
		saveFile(model.DictionaryInvert, fmt.Sprintf("%s/%s", outputFolder, "dictionaryInvert"))
		saveFile(model.Documents, fmt.Sprintf("%s/%s", outputFolder, "posting"))
		saveFile(model.ModelOptions, fmt.Sprintf("%s/%s", outputFolder, "modelOption"))
		end = time.Now()
		fmt.Println(end.Sub(start), "to save files")

	},
}

func saveFile(data interface{}, fileName string) {

	o, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(o)
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(inv)

	// Here you will define your flags and configuration settings.

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	inv.Flags().StringVarP(&inputFile, "file", "f", "", "The location of CACM collection (required)")
	inv.Flags().StringVarP(&outputFolder, "output-folder", "o",
		"", "The location to output the final files (required)")
	inv.MarkFlagRequired("file")
	inv.MarkFlagRequired("output-folder")
	inv.Flags().StringVarP(&stopLimit, "stop-word", "s", "", "add stop word removal")
	inv.Flags().BoolVarP(&porter, "porter", "p", false, "enable Porter's Stemming algorithm")
	//inv.Flags().BoolVarP(&lower, "lower-case", "l", false, "lower case each term")

}

//Load the files and the return is the a list of map which are the doc info
func loadFile(f string) []map[string]string {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(dat), "\n")
	lastToken := ""
	data := []map[string]string{}
	var d map[string]string
	var buffer bytes.Buffer

	for _, line := range lines {
		matches := re.FindAllString(line, -1)
		for _, match := range matches {
			m := match[1:]
			if strings.HasPrefix(m, "I") {
				if len(lastToken) > 0 {
					d[lastToken] = strings.TrimSpace(buffer.String())
				}
				if d != nil {
					data = append(data, d)
				}
				d = make(map[string]string)
				d["I"] = m[2:]
				buffer.Reset()
			} else {
				if len(lastToken) > 0 {
					d[lastToken] = strings.TrimSpace(buffer.String())
				}
				lastToken = m
				buffer.Reset()
			}
		}
		if matches == nil {
			buffer.Write([]byte(line))
			buffer.WriteRune('\n')
		}
	}
	if len(lastToken) > 0 {
		d[lastToken] = strings.TrimSpace(buffer.String())
	}
	if d != nil {
		data = append(data, d)
	}
	return data
}
