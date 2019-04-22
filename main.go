package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	defaultPath, err := defaultModfilePath()
	if err != nil {
		log.Fatal(err)
	}

	var path string
	flag.StringVar(&path, "modfile", defaultPath, "Path to go.mod")
	flag.Parse()

	err = parseModfile(path, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultModfilePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, "go.mod"), nil
}
