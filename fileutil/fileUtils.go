package fileutil

import (
	"io/ioutil"
	"fmt"
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