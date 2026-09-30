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

		log := LogEntry{Level: words[1], Endpoint: words[2]}
		return log, true
	}

	// or "" and false if the line is too short

	return LogEntry{}, false
}

func processFile(path string, level string) (map[string]int, map[string]int) {

	endpointCount := make(map[string]int)
	levelsCount := make(map[string]int)

	// Job 2: open the file.
	// If opening fails, print the error and exit.
	// Then make sure the file gets closed.

	file, err := os.Open(path)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer file.Close()

	// Job 3: create a scanner over the file and loop
	// for as long as there is another line.
	// Job 4: inside the loop, print the line.

	scanner := bufio.NewScanner(file)

	const maxCapacity = 1024 * 1024
	buf := make([]byte, 0, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		entry, ok := parseLine(scanner.Text())

		if !ok {
			continue
		}

		if level == "" || entry.Level == level {
			endpointCount[entry.Endpoint]++
			levelsCount[entry.Level]++
			/* fmt.Println(entry.Endpoint) */
		}

	}

	// After the loop: check whether the scanner hit an error.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return endpointCount, levelsCount
}

func main() {
	// Job 1: define a string flag called "file", then parse flags.

	// If the user gave no filename, print a message and exit.

	files := flag.String("files", "", "comma seperated files")

	

	l := flag.String("level", "", "Log Level")

	flag.Parse()

	paths := strings.Split(*files, ",")

	if *files == "" {
		fmt.Fprintln(os.Stderr, errors.New("Empty value"))
		os.Exit(1)
	}

	combinedEndpoints := make(map[string]int)
	combinedLevels := make(map[string]int)

	for _, path := range paths {
    endpoints, levels := processFile(path, *l)
    // merge endpoints into combinedEndpoints

	for key, value := range endpoints {
        combinedEndpoints[key] += value
    }
    // merge levels into combinedLevels

	for key, value := range levels {
        combinedLevels[key] += value
    }
}
	for key, value := range combinedEndpoints {
		fmt.Println(key, value)
	}
	for key, value := range combinedLevels {
		fmt.Println(key, value)
	}

}
