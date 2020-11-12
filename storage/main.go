package storage

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tusd "github.com/tus/tusd/pkg/handler"
)

var (
	uploadDir  string
	nameLength *int = flag.Int("file-length", 5, "Length of generated filename, extension omitted")
)

func LoadConfig() {
	uploadDir = flag.Args()[0]
}

func ProcessFile(event tusd.HookEvent) {
	removeInfoFile(event.Upload.ID)
}

func removeInfoFile(name string) {
	path := filepath.Join(uploadDir, fmt.Sprintf("%s.info", name))
	os.Remove(path)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
