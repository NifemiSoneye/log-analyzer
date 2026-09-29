package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type LogEntry struct {
	Level    string
	Endpoint string
}

func parseLine(line string) (LogEntry, bool) {

	// split, check length, return the level and true,
	words := strings.Fields(line)

	if len(words) >= 4 {

		log := LogEntry{Level :words[1], Endpoint:words[2]}
		return log, true
	}

	// or "" and false if the line is too short

	return LogEntry{}, false
}

func main() {
	// Job 1: define a string flag called "file", then parse flags.

	// If the user gave no filename, print a message and exit.

	f := flag.String("file", "", "parsed file")

	l := flag.String("level" , "" , "Log Level")

	flag.Parse()

	if *f == "" {
		fmt.Fprintln(os.Stderr, errors.New("Empty value"))
		os.Exit(1)
	}

	endpointCount := make(map[string]int)
	levelsCount := make(map[string]int)

	// Job 2: open the file.
	// If opening fails, print the error and exit.
	// Then make sure the file gets closed.

	file, err := os.Open(*f)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer file.Close()

	// Job 3: create a scanner over the file and loop
	// for as long as there is another line.
	// Job 4: inside the loop, print the line.

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, ok := parseLine(scanner.Text())

		if !ok {
			continue
		}

		
		

		if *l == "" || entry.Level == *l {
			endpointCount[entry.Endpoint]++
			levelsCount[entry.Level]++
			fmt.Println( entry.Endpoint )
		}
		
	}
	

	// After the loop: check whether the scanner hit an error.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for key , value := range endpointCount {
			fmt.Println(key , value)
		}
	for key , value := range levelsCount {
			fmt.Println(key , value)
		}
}
