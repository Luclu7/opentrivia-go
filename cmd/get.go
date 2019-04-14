package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	resty "gopkg.in/resty.v1"
	"html"
	"log"
	"math/rand"
	"time"
)

// cfe panic in case of an error
func cfe(err error) bool {
	if err != nil {
		log.Panicln(err)
		return false
	}
	return true
}

type Results struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

type Question struct {
	Results []Results
}

func fetchQuestion(nomber_of_questions_to_get string, type_of_question string) (Question, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Get("https://opentdb.com/api.php?amount=" + nomber_of_questions_to_get + "&type=" + type_of_question)
	cfe(err)
	var f Question
	err = json.Unmarshal(resp.Body(), &f)
	return f, err
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a question from OpenTriviaDB",
	Long:  "Get a question and its answers, shuffle them to get a random order and print them.",
	Run: func(cmd *cobra.Command, args []string) {
		question, err := fetchQuestion("1", "multiple")
		//	fmt.Println(question.Results)
		cfe(err)
		for _, results := range question.Results {
			println("Category: " + html.UnescapeString(results.Category))
			println("Type: " + html.UnescapeString(results.Type))
			println("Difficulty: " + html.UnescapeString(results.Difficulty))
			println("Question: " + html.UnescapeString(results.Question))
			println("Correct answer: " + html.UnescapeString(results.CorrectAnswer))
			println("Incorrect answers: " + html.UnescapeString(results.IncorrectAnswers[0]) + ", " + html.UnescapeString(results.IncorrectAnswers[1]) + ", " + html.UnescapeString(results.IncorrectAnswers[2]))
			var a [4]string
			a[0] = html.UnescapeString(results.CorrectAnswer)
			a[1] = html.UnescapeString(results.IncorrectAnswers[0])
			a[2] = html.UnescapeString(results.IncorrectAnswers[1])
			a[3] = html.UnescapeString(results.IncorrectAnswers[2])
			println("\nPossible answers:")
			rand.Seed(time.Now().UnixNano())
			// every examples were with ints, but it also works with strings, so ¯\_(ツ)_/¯
			rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
			println("A: " + a[0])
			println("D: " + a[1])
			println("C: " + a[2])
			println("D: " + a[3])
		}
		cfe(err)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("full", "f", false, "Help message for toggle")

}
