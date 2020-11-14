package storage

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
	"tus-server/config"
)

var (
	nameLength int
	chars      []rune
)

func init() {
	chars = generatePossibleChars()
}

func ProcessFile(file string) (newPath string) {
	removeInfoFile(file)
	name := generateNameBis(config.UploadDir, ".jpg")
	newPath = filepath.Join(config.UploadDir, name)
	err := os.Rename(filepath.Join(config.UploadDir, file), newPath)

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

func removeInfoFile(name string) {
	path := filepath.Join(config.UploadDir, fmt.Sprintf("%s.info", name))
	os.Remove(path)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
