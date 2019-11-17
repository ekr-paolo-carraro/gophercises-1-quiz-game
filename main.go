package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type questAnswer struct {
	question string
	answer   string
}

type inputQuest struct {
	inputReader bufio.Reader
}

func (iq inputQuest) askQuestion(question string, correctResult string) (bool, error) {
	fmt.Printf("Insert result for \"%v\": ", question)
	inputed, err := iq.inputReader.ReadString('\n')
	if err != nil {
		return false, err
	}

	inputed = strings.TrimSpace(inputed)
	if strings.ToLower(inputed) != strings.ToLower(correctResult) {
		return false, nil
	}

	return true, nil
}

var stopTime int
var questSource string

func main() {

	flag.IntVar(&stopTime, "time", 10, "duration time to answer quiz questions")
	flag.IntVar(&stopTime, "t", 10, "duration time to answer quiz questions (shorthand)")
	flag.StringVar(&questSource, "source", "problems.csv", "csv source for quiz questions in the form question, result")
	flag.StringVar(&questSource, "s", "problems.csv", "csv source for quiz questions in the form question, result (shorthand)")
	flag.Parse()

	var err error
	questions, err := loadQuestions(questSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	inputQuest := inputQuest{
		inputReader: *bufio.NewReader(os.Stdin),
	}

	fmt.Println("Start with quiz!")
	go timing(stopTime)

	var correctCounter int = 0
	var wrongQuestions []questAnswer = []questAnswer{}
	for i := 0; i < len(questions); i++ {
		var tempQuest questAnswer = questions[i]
		result, err := inputQuest.askQuestion(tempQuest.question, tempQuest.answer)
		if err != nil {
			log.Fatal(fmt.Sprintf("error on manage input: %v \n", err.Error()))
		}
		if result == true {
			correctCounter++
		} else {
			wrongQuestions = append(wrongQuestions, tempQuest)
		}
	}

	fmt.Println("-----------------------")
	fmt.Printf("Correct results %v\n", correctCounter)
	fmt.Println("-----------------------")
	if len(wrongQuestions) > 0 {
		fmt.Println("Wrong quests:")
		for i := 0; i < len(wrongQuestions); i++ {
			fmt.Printf("- %v. Correct response: %v\n", wrongQuestions[i].question, wrongQuestions[i].answer)
		}
	}

}

func loadQuestions(source string) ([]questAnswer, error) {
	questionsFile, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("Error on open resource file: %v", err.Error())
	}
	csvReader := csv.NewReader(questionsFile)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error on reading resource file: %v", err.Error())
	}

	var results []questAnswer = make([]questAnswer, len(csvData))
	for i, questItem := range csvData {
		results[i] = questAnswer{
			question: questItem[0],
			answer:   questItem[1],
		}
	}

	return results, nil
}

func timing(stopTime int) {
	time.Sleep(time.Duration(time.Duration(stopTime) * time.Second))
	fmt.Println("")
	fmt.Println("Time is over!")
	os.Exit(1)
}
