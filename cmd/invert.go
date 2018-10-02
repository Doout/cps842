package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/doout/cps842/pkg/document"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
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
		data := loadFile(inputFile)
		var doc *document.Documents
		tokenParser := []func(token string) string{}
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
		if porter {
			tokenParser = append(tokenParser, document.PorterStemmer)
		}
		if lower {
			tokenParser = append(tokenParser, document.ToLower)
		}
		if tokenParser != nil {
			doc = document.BuildDocumentWithTokenParser(data, tokenParser...)
		} else {
			doc = document.BuildDocument(data)
		}
		_ = doc
		//save dictionary
		fmt.Println(fmt.Sprintf("Saving files here %s", outputFolder))
		saveFile(doc.GetDictionarySort(), fmt.Sprintf("%s/%s", outputFolder, "dictionary"))
		saveFile(doc, fmt.Sprintf("%s/%s", outputFolder, "postings"))
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
	inv.Flags().BoolVarP(&lower, "lower-case", "l", false, "lower case each term")

}

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
	return data
}
