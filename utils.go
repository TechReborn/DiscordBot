package main

import (
	"bufio"
	"bytes"
	"net/http"
	"os"
	"time"
	"unsafe"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string) (string, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	var buf = new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var b = buf.Bytes()
	var s = *(*string)(unsafe.Pointer(&b))
	return s, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
