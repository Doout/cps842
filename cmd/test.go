package cmd

import (
	"github.com/spf13/cobra"

	"bufio"
	"encoding/json"
	"fmt"
	"github.com/doout/cps842/pkg/document"
	"io/ioutil"
	"os"
	"strings"
)

// inv represents the playbook command
var test = &cobra.Command{
	Use:   "test",
	Short: "Test inverted index",
	Long:  `Take a collection of documents and generate its inverted index.`,
	Run: func(cmd *cobra.Command, args []string) {
		doc := document.Documents{}
		loadJsonFromFile(&doc, "./data/postings")
		doc.GetFristDocSum("writable")
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter text: ")
			input, _ := reader.ReadString('\n')
			input = strings.Trim(input, "\n")
			if input == "ZZEND" {
				break
			}
			fmt.Println(doc.GetFristDocSum(input))
		}

	},
}

func loadJsonFromFile(t interface{}, file string) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(dat, t); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(test)

}
