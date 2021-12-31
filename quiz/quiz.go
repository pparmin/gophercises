package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
	"unicode"
)

type problem struct {
	question string
	answer   string
}

func parse(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	// filter all non-number input (including extra space)
	input = strings.TrimLeftFunc(input, func(r rune) bool {
		return !unicode.IsNumber(r)
	})
	input = strings.TrimRightFunc(input, func(r rune) bool {
		return !unicode.IsNumber(r)
	})
	return input
}

func readFile(path string) []problem {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	quiz := csv.NewReader(file)
	questions, err := quiz.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	problems := parse(questions)
	return problems
}

func exit(correct, total int) {
	fmt.Printf("\nCongratulations! You answered %d out of %d questions correctly!\n", correct, total)
	os.Exit(0)
}

func play(problems []problem, limit int) {
	correct := 0
	// i := 1
	start := time.Now()
	go func() {
		for {
			if math.Floor(time.Since(start).Seconds()) == float64(limit) {
				exit(correct, len(problems))
			}
		}
	}()

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)
		input := getInput()
		if input == p.answer {
			fmt.Println("Correct!")
			correct++
		}
		if i == len(problems)-1 {
			exit(correct, len(problems))
		}
	}
}

func main() {
	csvFile := flag.String("csv", "csv/default.csv", "a csv file in the format of 'question, answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	problems := readFile(*csvFile)
	play(problems, *limit)
}
