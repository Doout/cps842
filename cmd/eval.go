package cmd

import (
	"fmt"
	"github.com/doout/cps842/pkg/document"
	"github.com/doout/cps842/pkg/pagerank"
	"github.com/doout/cps842/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// inv represents the playbook command
var eval = &cobra.Command{
	Use:   "eval",
	Short: "Eval the search",
	Long:  `Evaluate the performance of the IR system, output will be the average MAP and R-Precision values over all queries`,
	Run: func(cmd *cobra.Command, args []string) {
		timer := utils.Timer{}
		timer.Start()
		model := document.LoadModel(folder)
		timer.Stop()
		fmt.Println(timer.Duration(), "to load the model")
		timer.Start()

		//Load in the pageRank

		pageRankGraph := pagerank.New()
		utils.LoadJsonFromFile(&pageRankGraph, fmt.Sprintf("%s/%s", folder, "pagerank_graph"))
		prs := pageRankGraph.Rank(0.85, 1e-10)
		prs.Norm()
		timer.Stop()
		fmt.Println(timer.Duration(), "to load/compute page rank")
		//qrels.text output
		// index doc ID
		//01 1410  0 0
		qrels := LoadQrels(qrels)
		_ = qrels
		querys := loadFile(query)
		re = regexp.MustCompile(`(?m)^[ \d]+.`)
		//clean up N
		for index := range querys {
			m := re.FindString(querys[index]["N"])
			if strings.HasPrefix(querys[index]["N"], m) {
				querys[index]["N"] = strings.Trim(string([]byte(querys[index]["N"])[len(m):]), " ")
			}
		}

		matchs := make(map[int][]document.Result)
		for _, q := range querys {
			//We only support W
			qId, err := strconv.Atoi(q["I"])
			if err != nil {
				panic(err)
			}
			if qId == 0 {
				break
			}
			matchs[qId], _ = model.SearchWithPageRank(q, 0.7, 0.3, prs)
		}

		keys := []int{}
		for key, _ := range qrels {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		totalAP := float64(0)
		fmt.Println("R-Precision values")
		for _, key := range keys {
			ap, rp := GetAPAndRP(qrels[key], matchs[key])
			fmt.Printf("Query: %d\tAP: %0.2f\tR-Precision: %0.2f\n", key, ap, rp)
			totalAP += ap
		}
		MAP := totalAP / float64(len(keys))
		fmt.Println("MAP:", MAP)
		timer.Stop()

		fmt.Println(timer.Duration(), "search/compute R-Precision")
		_ = model

	},
}

func GetAPAndRP(docSet []int, results []document.Result) (float64, float64) {
	union, ranks := IntersectionWithIndex(docSet, results)
	total := float64(0)
	for index, _ := range union {
		total += (float64(index+1) / float64(ranks[index]))
	}
	return total / float64(len(docSet)), float64(float64(len(union)) / float64(len(results)))
}

func IntersectionWithIndex(docSet []int, results []document.Result) ([]int, []int) {
	union := []int{}
	indexs := []int{}
	for index, result := range results {
		for _, docId := range docSet {
			if docId == int(result.Document) {
				union = append(union, docId)
				indexs = append(indexs, index+1)
			}
		}
	}
	return union, indexs
}

func LoadQrels(f string) map[int][]int {
	qrels := make(map[int][]int)
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(dat), "\n")
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		if len(tokens) < 2 {
			continue
		}
		index, err := strconv.Atoi(tokens[0])
		if err != nil {
			panic(err)
		}
		docID, err := strconv.Atoi(tokens[1])
		if err != nil {
			panic(err)
		}
		if _, ok := qrels[index]; ok {
			qrels[index] = append(qrels[index], docID)
		} else {
			qrels[index] = []int{docID}
		}

	}
	return qrels
}

var (
	qrels string
	query string
)

func init() {
	rootCmd.AddCommand(eval)
	eval.Flags().StringVarP(&folder, "folder", "f", "", "Folder location where the posting/doc files are (required)")
	eval.MarkFlagRequired("folder")

	eval.Flags().StringVarP(&qrels, "qrels", "r", "", "The qrels.text file")
	eval.Flags().StringVarP(&query, "query", "q", "", "The query file")
}
