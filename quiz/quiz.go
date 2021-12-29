package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"
	"unicode"
)

func exit(correct, total int) {
	fmt.Printf("\nCongratulations! You answered %d out of %d questions correctly!\n", correct, total)
	os.Exit(0)
}

func play(quiz *csv.Reader, limit int, total int) {
	correct := 0
	i := 1
	start := time.Now()
	go func() {
		for {
			if math.Floor(time.Since(start).Seconds()) == float64(limit) {
				exit(correct, total)
			}
		}
	}()

	for {
		record, err := quiz.Read()
		if err == io.EOF {
			if err != nil {
				log.Fatal(err)
			}
			exit(correct, total)
		}

		if err != nil {
			log.Fatal("PLAY FUNCTION ", err)
		}
		question := record[0]
		solution := record[1]
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Problem #%d: %s = ", i, question)
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

		if input == solution {
			fmt.Println("Correct!")
			correct++
		}

		i++
	}
}

func main() {
	csvFile := flag.String("csv", "csv/default.csv", "a csv file in the format of 'question, answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// second file Reader in order to get total number of questions
	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	f, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}

	quiz := csv.NewReader(file)
	q := csv.NewReader(f)

	questions, err := q.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	total := len(questions)
	play(quiz, *limit, total)
}
