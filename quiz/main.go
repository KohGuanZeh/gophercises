package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	filename, timer, shuffle := parseFlags()
	problems := getProblems(filename)
	problems = checkProblems(problems)
	if shuffle {
		shuffleProblems(problems)
	}

	fmt.Println("You have", timer, "seconds to answer", len(problems), "questions.")
	fmt.Println("Press enter to start.")
	_, err := fmt.Scanln()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	timerChannel := make(chan bool)
	go runTimer(timerChannel, timer)
	runQuiz(timerChannel, problems)
}

func parseFlags() (string, uint, bool) {
	filenamePtr := flag.String("file", "problems.csv", "CSV file name containing problem set and answers. Default is problems.csv.")
	flag.StringVar(filenamePtr, "f", *filenamePtr, "Alias for --file.")
	timerPtr := flag.Uint("time", 30, "Time in seconds to answer all questions. Default is 30s.")
	flag.UintVar(timerPtr, "t", *timerPtr, "Alias for --time.")
	shufflePtr := flag.Bool("shuffle", false, "Shuffle question order. Default is false.")
	flag.BoolVar(shufflePtr, "s", *shufflePtr, "Alias for --shuffle.")
	flag.Parse()

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		log.Fatalln("Failed to parse command-line flags: ", err)
	}

	return *filenamePtr, *timerPtr, *shufflePtr
}

func getProblems(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to read file %s: %s\n", filename, err)
	}

	r := csv.NewReader(f)
	problems, err := r.ReadAll()
	f.Close()
	if err != nil {
		log.Fatalln("Error reading problems:", err)
	}
	return problems
}

func checkProblems(problems [][]string) [][]string {
	checkedProblems := make([][]string, len(problems))
	i := 0
	for _, problem := range problems {
		if len(problem) < 2 {
			log.Println("Error: Question or answer seems to be missing.")
			log.Println("Skipping.")
			continue
		}
		if len(problem) > 2 {
			log.Println("Warning: More items than expected from problem.")
			log.Println("Only accepting the first 2 values.")
		}
		cleanedQns := strings.TrimSpace(problem[0])
		cleanedAns := strings.ToLower(strings.TrimSpace(problem[1]))
		checkedProblem := []string{cleanedQns, cleanedAns}
		checkedProblems[i] = checkedProblem
		i += 1
	}
	return checkedProblems
}

func shuffleProblems(problems [][]string) {
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

func runTimer(timerChannel chan bool, timer uint) {
	time.Sleep(time.Duration(timer) * time.Second)
	timerChannel <- true
}

func runQuiz(timerChannel chan bool, problems [][]string) {
	userAnsChannel := make(chan bool)
	total, correct := len(problems), 0
	stop := false
	for i, problem := range problems {
		fmt.Println("Question ", i+1)
		qns, ans := problem[0], problem[1]
		fmt.Println(qns)
		var userAns string
		go getUserInput(userAnsChannel, &userAns)
		select {
		case <-timerChannel:
			fmt.Println("Time's up!")
			stop = true
		case <-userAnsChannel:
			userAns = strings.ToLower(strings.TrimSpace(userAns))
			if userAns == ans {
				correct += 1
				fmt.Print("Correct! ")
			} else {
				fmt.Println("Wrong! ")
			}
			fmt.Printf("Score: %d/%d\n", correct, total)
			fmt.Println()
		}
		if stop {
			break
		}
	}
	fmt.Printf("Congratulations, you scored a total of %d/%d!\n", correct, total)
}

func getUserInput(userAnsChannel chan bool, userAns *string) {
	n, err := fmt.Scanln(userAns)
	if err != nil && n != 0 {
		log.Fatalln("Error:", err)
	} else if n == 0 {
		*userAns = ""
	}
	userAnsChannel <- true
}
