package fileScrape

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func loadFile(fileName string) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new Scanner for the file
	scanner := bufio.NewScanner(file)

	// Loop over all lines in the file and print them
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Check for errors during Scan
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
