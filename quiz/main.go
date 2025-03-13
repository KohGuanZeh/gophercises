package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main() {
	filenamePtr := flag.String("file", "problems.csv", "CSV file name containing problem set and answers.")
	flag.StringVar(filenamePtr, "f", *filenamePtr, "Alias for --file.")
	shufflePtr := flag.Bool("shuffle", false, "Shuffle question order.")
	flag.BoolVar(shufflePtr, "s", *shufflePtr, "Alias for --shuffle.")
	flag.Parse()

	f, err := os.Open(*filenamePtr)
	if err != nil {
		log.Fatalln("Unable to read file", *filenamePtr, ":", err)
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	f.Close()
	if err != nil {
		log.Fatalln("Error reading records:", err)
	}

	if *shufflePtr {
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}
	total, correct := 0, 0
	for _, record := range records {
		if len(record) < 2 {
			log.Println("Error: Question or answer seems to be missing.")
			log.Println("Skipping.")
		}
		total += 1
		fmt.Println()
		fmt.Println("Question ", total)
		qns := strings.TrimSpace(record[0])
		fmt.Println(qns)
		var userAns string
		_, err := fmt.Scanln(&userAns)
		if err != nil {
			log.Fatalln("Error:", err)
		}
		userAns = strings.TrimSpace(userAns)
		ans := strings.TrimSpace(record[1])
		if userAns == ans {
			correct += 1
			fmt.Print("Correct! ")
		} else {
			fmt.Println("Wrong! ")
		}
		fmt.Println("Score:", correct, "/", total)
	}

	fmt.Println()
	fmt.Println("Congratulations, you scored a total of", correct, "/", total, "!")
}
