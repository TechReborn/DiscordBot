package file

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//ReadString reads a string from a file
func ReadString(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//WriteString writes a string to a file
func WriteString(str string, file string) error {
	return ioutil.WriteFile(file, []byte(str), 0644)
}

//AppendString appends a string to a file, or creates a new file with the string if the file does not exist
func AppendString(str string, file string) error {
	if Exists(file) {
		fileStr, err := ReadString(file)
		if err != nil {
			return err
		}
		return WriteString(fileStr+"\n"+str, file)
	} else {
		return WriteString(str, file)
	}
}

//Exists checks to see if a file exists
func Exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

//ReadLines reads each line of the file into a string array
func ReadLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func MakeDir(fileName string) error {
	return os.MkdirAll(fileName, os.ModePerm)
}

func GetRunPath() (string, error) {
	ex, err := os.Getwd()
	if err != nil {
		return "", err
	}
	exPath := path.Dir(ex)
	return exPath, nil
}

func DeleteDir(dir string) error {
	if !Exists(dir) {
		return errors.New("File not found")
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func FormatPath(path string) string {
	return strings.Replace(path, "/", string(os.PathSeparator), -1)
}
