package config

import (
	"flag"

	_ "github.com/joho/godotenv/autoload"
)

func Load() {
	flag.Parse()
}
