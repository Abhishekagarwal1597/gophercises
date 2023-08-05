/* Create a program that will read in a quiz provided via a CSV file (more details below) and will then give the quiz to a user keeping track of how many questions they get right and how many they get incorrect. Regardless of whether the answer is correct or wrong the next question should be asked immediately afterwards.

The CSV file should default to problems.csv (example shown below), but the user should be able to customize the filename via a flag.

The CSV file will be in a format like below, where the first column is a question and the second column in the same row is the answer to that question.*/

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	// Open a CSV file
	fh, err := os.Open("problem.csv")
	if err != nil {
		fmt.Errorf("Error Reading CSV", err.Error())
	}
	defer fh.Close()
	// Get a file Reader
	fileReader := csv.NewReader(fh)
	// Read all the file content
	data, err := fileReader.ReadAll()
	if err != nil {
		fmt.Errorf("Error Reading CSV", err.Error())
	}

	// create a channel for taking input data
	input_chan := make(chan string)
	var pos, neg int

	for key, value := range data {

		fmt.Printf("Your %d question:\n", key+1)
		var user_input int

		go func() {
			fmt.Scan(&user_input)
			input_chan <- "data input"
		}()

		select {
		case <-input_chan:
			// fmt.Println("data received")
			result := CompareData(value[0], user_input)
			if result {
				pos++
			} else {
				neg++
			}
		case <-terminateSignal():
			fmt.Println("channel terminated")
			os.Exit(1)
		}

	}
	fmt.Printf("result: Correct- : %d, Wrong: %d", pos, neg)

}

// terminate signal if user doesn't enter any data
func terminateSignal() chan os.Signal {
	terminateCh := make(chan os.Signal, 1)
	signal.Notify(terminateCh, syscall.SIGINT, syscall.SIGTERM)
	return terminateCh
}

// compare userInput with the desired output
func CompareData(InputData string, userInput int) bool {
	split := strings.Split(InputData, "+")
	sum := 0
	for _, val := range split {
		num, _ := strconv.Atoi(val)
		sum += num
	}

	if userInput == sum {
		return true
	}
	return false
}
