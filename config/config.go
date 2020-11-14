package config

import (
	"flag"

	_ "github.com/joho/godotenv/autoload"
)

var (
	UploadDir      string
	FileNameLength int
)

func Load() {
	flag.IntVar(&FileNameLength, "file-length", 5, "Length of generated filename, extension omitted")
	flag.Parse()
	UploadDir = flag.Args()[0]
}
