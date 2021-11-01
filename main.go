package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	LOGO = `
_     _       __ 
| |   (_) ___ / _|
| |   | |/ __| |_ 
| |___| | (__|  _|
|_____|_|\___|_|  


`
	FILE_CONFIG = "config.yaml"
)

func main() {

	if GetConfig().Logo.Show {
		color.Yellow(LOGO)
	}

	path := flag.String("path", "", "Specifies path of JSON file.")
	pattern := flag.String("pattern", "", "Search pattern")

	flag.Parse()

	iterateFiles(*path, *pattern)

}

// Iterate list files
func iterateFiles(directory, pattern string) {

	if directory == "" {
		directory = getCurrentDirectory()
	}

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		// Ignore direcotry or Files defined in config.yaml
		if GetConfig().File.CheckIgnore(file.Name()) {
			continue
		}

		ext := strings.Replace(filepath.Ext(file.Name()), ".", "", -1)

		filePath := directory + string(os.PathSeparator) + file.Name()

		if GetConfig().File.CheckFormat(ext) && (ext == "json" || ext == "yaml" || ext == "yml") {
			line := Reader{filePath}.Search(pattern)
			yellow := color.New(color.FgYellow).SprintFunc()

			fmt.Print(yellow(filePath + ":" + strconv.Itoa(line) + " : "))

			if line > -1 {
				color.Green("True")
			} else {
				color.Red("False")
			}

		} else if file.IsDir() && GetConfig().File.Recursion {
			iterateFiles(filePath, pattern)
		}

	}
}

// DIRECTORY
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
