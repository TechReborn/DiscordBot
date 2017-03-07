package fileutil

import (
	"io/ioutil"
	"fmt"
	"os"
	"log"
	"bufio"
)

func ReadStringFromFile(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

func WriteStringToFile(str string, file string) {
	ioutil.WriteFile(file, []byte(str), 0644)
}

func AppendStringToFile(str string, file string) {
		if FileExists(file) {
				WriteStringToFile( ReadStringFromFile(file) + "\n" + str, file)
		} else {
				WriteStringToFile(str, file)
		}
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}


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