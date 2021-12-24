package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func play(quiz *csv.Reader) {
	counter := 0
	i := 1
	for {
		record, err := quiz.Read()
		if err == io.EOF {
			fmt.Printf("Congratulations! You answered %d out of %d questions correctly!\n", counter, i-1)
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		question := record[0]
		solution := record[1]
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Problem #%d: %s = ", i, question)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		input = strings.TrimRight(input, "\n")

		if input == solution {
			fmt.Println("Correct!")
			counter++
		}
		i++
	}
}

func main() {
	csvFile := flag.String("csv", "default.csv", "a csv file in the format of 'question, answer' (default \"problems.csv\")")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds (default 30)")
	flag.Parse()

	fmt.Println("LIMIT: ", *timeLimit)

	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	quiz := csv.NewReader(file)
	play(quiz)
}
