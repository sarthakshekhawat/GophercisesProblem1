package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problemSet struct {
	question string
	answer   string
}

// Read the csv file and stores the data in "quiz"
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
	correctAnswers := 0 // keeps the track of total number of correct answers
	if *shuffleQuestions {
		shuffle()
	}
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
			response = whiteSpaceRemover(response)     // removes the leading and trailing white space
			if strings.EqualFold(response, q.answer) { // compares the string in a case in-sensitive manner
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

// Removes the white space from the string
func whiteSpaceRemover(response string) string {
	return strings.TrimSpace(response)
}

func result(correctAnswers int) {
	fmt.Printf("Out of %v questions, %v questions are correct \n", len(quiz), correctAnswers)
}

//Shuffle the questions in the quiz
func shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(quiz), func(i, j int) { quiz[i], quiz[j] = quiz[j], quiz[i] })
}

//This slice will store all the questions and answers of the quiz
var quiz []problemSet

//Flag, for the csv loction
var csvFileLocation = flag.String("location", "problems.csv", "Location of the csv")

//Flag, for quiz timer
var quizTimer = flag.Int64("timer", 30, "Duration of the quiz, in seconds")

//Flag, to shuffle the questions
var shuffleQuestions = flag.Bool("shuffle", true, "Shuffle the questions in the quiz")

func main() {

	flag.Parse()
	csvReader()
	startQuiz()
}
