package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "questions.csv", "CSV file with questions and answers")
	timeLimit := flag.Int("limit", 30, "Time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open the CSV file: %s\n", *csvFilename)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Failed to parse the CSV file.")
		os.Exit(1)
	}

	questions := parseQuestions(records)

	fmt.Println("Press Enter to start the quiz.")
	fmt.Scanln()

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	questionCh := make(chan bool)

	go func() {
		for i, q := range questions {
			fmt.Printf("Question %d: %s = ", i+1, q.question)
			var answer string
			fmt.Scanln(&answer)

			if strings.TrimSpace(answer) == q.answer {
				correct++
			}
			questionCh <- true
		}
		close(questionCh)
	}()

	for {
		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			printResults(correct, len(questions))
			return
		case _, ok := <-questionCh:
			if !ok {
				printResults(correct, len(questions))
				return
			}
		}
	}
}

type Question struct {
	question string
	answer   string
}

func parseQuestions(records [][]string) []Question {
	questions := make([]Question, len(records))
	for i, record := range records {
		questions[i] = Question{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}
	return questions
}

func printResults(correct, total int) {
	fmt.Printf("\nYou got %d out of %d questions correct.\n", correct, total)
}
