package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
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

func processFile(path string, level string, combinedEndpoints map[string]int, combinedLevels map[string]int, mu *sync.Mutex, wg *sync.WaitGroup) {

	defer wg.Done()
	endpointCount := make(map[string]int)
	levelsCount := make(map[string]int)

	file, err := os.Open(path)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer file.Close()

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

	mu.Lock()
for key, value := range endpointCount {
    combinedEndpoints[key] += value
}
for key, value := range levelsCount {
    combinedLevels[key] += value
}
mu.Unlock()
}

func main() {
	// Job 1: define a string flag called "file", then parse flags.

	// If the user gave no filename, print a message and exit.

	files := flag.String("files", "", "comma seperated files")

	

	l := flag.String("level", "", "Log Level")

	flag.Parse()

	paths := strings.Split(*files, ",")

	if *files == "" {
		fmt.Fprintln(os.Stderr, "usage: loganalyzer -files=<path1,path2,...> [-level=LEVEL]")
		os.Exit(1)
	}

	var mu sync.Mutex
var wg sync.WaitGroup
	combinedEndpoints := make(map[string]int)
	combinedLevels := make(map[string]int)

	for _, path := range paths {
	wg.Add(1)
    go processFile(path, *l , combinedEndpoints , combinedLevels , &mu , &wg)
}
	wg.Wait()

	if _, ok := combinedLevels[*l] ; !ok  && *l != "" {
		fmt.Fprintln(os.Stderr, "no lines matched level" , *l)
	}
	for key, value := range combinedEndpoints {
		fmt.Println(key, value)
	}
	for key, value := range combinedLevels {
		fmt.Println(key, value)
	}

}
