package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// get quiz game options
	var csvFilename = flag.String("csv", "problems.csv", "quiz game problems file")
	var timeLimit = flag.Int("limit", 30, "timeLimit for the quiz before it expires")
	flag.Parse()

	// open file
	f, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Could not open csv file, %s. %s", *csvFilename, err)
		os.Exit(1)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	correct_answers := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answers := make(chan string)
	err = nil

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	for err != io.EOF {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s:", rec[0])
		go func() {
			var answer string
			_, err := fmt.Scan(&answer)
			if err != nil {
				log.Fatal(err)
				return
			}
			answers <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime expired")
			fmt.Printf("You scored %d!", correct_answers)
			return
		case val := <-answers:
			if val == rec[1] {
				correct_answers++
			}
		}
	}
	fmt.Printf("You scored %d!", correct_answers)
}
