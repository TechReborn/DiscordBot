package fileutil

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//ReadStringFromFile reads a string from a file
func ReadStringFromFile(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

//WriteStringToFile writes a string to a file
func WriteStringToFile(str string, file string) {
	ioutil.WriteFile(file, []byte(str), 0644)
}

//AppendStringToFile appends a string to a file, or creates a new file with the string if the file does not exist
func AppendStringToFile(str string, file string) {
	if FileExists(file) {
		WriteStringToFile(ReadStringFromFile(file)+"\n"+str, file)
	} else {
		WriteStringToFile(str, file)
	}
}

//FileExists checks to see if a file exists
func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

//ReadLinesFromFile reads each line of the file into a string array
func ReadLinesFromFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}
