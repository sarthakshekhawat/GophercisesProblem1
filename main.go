package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problemSet struct {
	question string
	answer   string
}

func csvReader() {
	csvFile, err := os.Open(*csvFileLocation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, line := range csvLines {
		ques := problemSet{
			question: line[0],
			answer:   line[1],
		}
		quiz = append(quiz, ques)
	}
}

func startQuiz() {
	correctAnswers := 0
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Press \"Enter\" to start the Quiz!")
	scanner.Scan()

	timer := time.NewTimer(time.Duration(*quizTimer) * time.Second)

	for _, q := range quiz {
		fmt.Println(q.question)
		c := make(chan string)
		go func() {
			scanner.Scan()
			input := scanner.Text()
			c <- input
		}()

		select {
		case response := <-c:
			if response == q.answer {
				correctAnswers++
			}
		case <-timer.C:
			fmt.Println("TIME'S UP!")
			result(correctAnswers)
			return
		}
	}
	result(correctAnswers)
}

func result(correctAnswers int) {
	fmt.Printf("Out of %v questions, %v questions are correct \n", len(quiz), correctAnswers)
}

//This slice will store all the questions and answers of the quiz
var quiz []problemSet

//Flag, for the csv loction
var csvFileLocation = flag.String("location", "problems.csv", "Location of the file")

//Flag, for quiz timer
var quizTimer = flag.Int64("timer", 30, "Duration of the quiz, in seconds")

func main() {

	flag.Parse()
	csvReader()
	startQuiz()
}
