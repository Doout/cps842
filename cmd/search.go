package cmd

import (
	"bufio"
	"fmt"
	"github.com/bclicn/color"
	"github.com/doout/cps842/pkg/document"
	"github.com/doout/cps842/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

// inv represents the playbook command
var search = &cobra.Command{
	Use:   "search",
	Short: "Search the inverted index",
	Long:  `Search for the best documents match using cos-sim and output the best match`,
	Run: func(cmd *cobra.Command, args []string) {

		model := document.LoadModel(folder)
		if !Interact && Search == "" {
			fmt.Println("Nothing to search for.")
			os.Exit(1)
		}
		if !Interact {
			a := make(map[string]string)
			//We only looking at the body and nothing else.
			// model.Search does not support anything but W at this moment
			a["W"] = Search
			start := time.Now()
			res := model.Search(a)
			end := time.Now()
			fmt.Println("Doc ID\tCos-sim")
			for _, r := range res {
				fmt.Printf("%d\t%f\n", r.Document, r.Value)
			}
			fmt.Println(end.Sub(start), "to find the best match")
		} else {
			reader := bufio.NewReader(os.Stdin)
			timer := utils.Timer{}
			for {
				fmt.Print("Enter text: ")
				input, _ := reader.ReadString('\n')
				input = strings.Trim(input, "\n")
				if input == "ZZEND" {
					break
				}
				a := make(map[string]string)
				//We only looking at the body and nothing else.
				// model.Search does not support anything but W at this moment
				a["W"] = input
				timer.Start()
				res := model.Search(a)
				timer.Stop()
				for index, r := range res {
					docInfo := model.Info[int(r.Document)]
					//The output will be
					//Rank (index)
					//Title
					//	(Title)
					//Author names
					// (Author names)
					fmt.Printf("%s: %d\n%s:\n\t%s\n%s\n\t%s\n",
						color.BCyan("Rank"),
						index+1,
						color.BCyan("Title"),
						strings.Replace(docInfo["T"], "\n", " ", -1),
						color.BCyan("Author names:"),
						strings.Replace(docInfo["A"], "\n", " ", -1))
				}

				fmt.Println()
				fmt.Println(timer.Duration(), "to find the best match")
			}
		}
		_ = model

	},
}

var (
	Search   string
	Interact bool
)

func init() {
	// Don't load this sub command, for the project
	//rootCmd.AddCommand(search)
	search.Flags().StringVarP(&folder, "folder", "f", "", "Folder location where posting/doc files are (required)")
	search.MarkFlagRequired("folder")

	search.Flags().StringVarP(&Search, "search", "s", "", "What to search for")
	search.Flags().BoolVarP(&Interact, "interact-mode", "i", false, "Interact with the user")
}
