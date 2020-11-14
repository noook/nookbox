package storage

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	nameLength int
	chars      []rune
	uploadDir  string
)

func init() {
	flag.IntVar(&nameLength, "file-length", 5, "Length of generated filename, extension omitted")
	chars = generatePossibleChars()
}

func ProcessFile(file string) (newPath string) {
	removeInfoFile(uploadDir, file)
	name := generateNameBis(uploadDir, ".jpg")
	newPath = filepath.Join(uploadDir, name)
	err := os.Rename(filepath.Join(uploadDir, file), newPath)

	if err != nil {
		fmt.Println(err)
	}

	return name
}

func generatePossibleChars() (list []rune) {
	for i := 48; i <= 57; i++ {
		list = append(list, rune(i))
	}
	for i := 65; i <= 90; i++ {
		list = append(list, rune(i))
	}
	for i := 97; i <= 122; i++ {
		list = append(list, rune(i))
	}

	return list
}

func guid(length int) (identifier string) {
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < length; i++ {
		identifier += string(chars[rand.Intn(len(chars))])
	}

	return identifier
}

func generateName(path string, extension string) (string, error) {
	for i := 0; i < 10000; i++ {
		id := guid(nameLength)
		if !fileExists(filepath.Join(path, id, extension)) {
			return id + extension, nil
		}
	}

	return "", errors.New("Couldn't find a value withing 10 000 tries")
}

func generateNameBis(path string, extension string) string {
	for id := guid(nameLength); !fileExists(filepath.Join(path, id, extension)); id = guid(nameLength) {
		return id + extension
	}

	return guid(nameLength) + extension
}

func removeInfoFile(uploadDir, name string) {
	path := filepath.Join(uploadDir, fmt.Sprintf("%s.info", name))
	os.Remove(path)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
