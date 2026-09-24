package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Job 1: define a string flag called "file", then parse flags.

	// If the user gave no filename, print a message and exit.

	f := flag.String("file" , "" , "parsed file")

	flag.Parse()
	
	if *f == "" {
		fmt.Fprintln(os.Stderr, errors.New("Empty value"))
		os.Exit(1)
	}

	// Job 2: open the file.
	// If opening fails, print the error and exit.
	// Then make sure the file gets closed.


	file , err := os.Open(*f)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer file.Close()


	// Job 3: create a scanner over the file and loop
	// for as long as there is another line.
	// Job 4: inside the loop, print the line.

	scanner := bufio.NewScanner(file)

	for  scanner.Scan() {
		fmt.Println(scanner.Text())
	}


	

	// After the loop: check whether the scanner hit an error.
	 if err := scanner.Err() ; err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}