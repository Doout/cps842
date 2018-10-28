package cmd

import (
	"github.com/spf13/cobra"

	"encoding/json"

	"io/ioutil"
)

// inv represents the playbook command
var test = &cobra.Command{
	Use:   "test",
	Short: "Test inverted index",
	Long:  `Take the posting file generate from invert and test it.`,
	Run: func(cmd *cobra.Command, args []string) {
		//tokenParser := []func(token string) string{document.RemovePunctuation, document.ToLower}
		//_ = tokenParser
		//doc := document.TermFrequencys{}
		//loadJsonFromFile(&doc.TermFrequency, fmt.Sprintf("%s/%s", folder, "postings"))
		//loadJsonFromFile(&doc.Info, fmt.Sprintf("%s/%s", folder, "docinfo"))
		//reader := bufio.NewReader(os.Stdin)
		//totalTime := int64(0)
		//total := int64(0)
		//for {
		//	fmt.Print("Enter text: ")
		//	input, _ := reader.ReadString('\n')
		//	input = strings.Trim(input, "\n")
		//	if input == "ZZEND" {
		//		break
		//	}
		//	start := time.Now()
		//	output := doc.GetTermSum(input)
		//	if output == "" {
		//		continue
		//	}
		//	end := time.Now()
		//	totalTime += end.Sub(start).Nanoseconds()
		//	//We don't want to add the time it take to output to this as it not the lookup time
		//	fmt.Println(output)
		//	fmt.Println("Time: ", end.Sub(start))
		//	total++
		//
		//}
		//if total > 0 {
		//	fmt.Println("The average time is ", time.Duration(totalTime/total))
		//}

	},
}

var (
	folder string
)

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
	//Uncomment to add this subcommand back into the list
	//rootCmd.AddCommand(test)
	test.Flags().StringVarP(&folder, "folder", "f", "", "Folder location where posting/doc files are (required)")
	test.MarkFlagRequired("folder")
}
