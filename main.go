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

// Iterate all files from directory
func iterateFiles(directory, pattern string) {

	var totalFile int

	var recFunc func(dir string)

	recFunc = func(dir string) {

		if dir == "" {
			dir = getCurrentDirectory()
		}

		files, err := ioutil.ReadDir(dir)

		// if it is not a directory, checks if it is a file and append it to the files array
		// otherwise it returns an error
		if err != nil {
			f, err2 := os.Stat(dir)
			files = append(files, f)
			if err2 != nil || !f.Mode().IsRegular() {
				log.Fatal(err)
				log.Fatal(err2)
			}
		}

		for _, file := range files {

			// Ignore direcotry or Files defined in config.yaml
			if GetConfig().File.CheckIgnore(file.Name()) {
				continue
			}

			ext := strings.Replace(filepath.Ext(file.Name()), ".", "", -1)

			filePath := dir

			if !strings.Contains(dir, file.Name()) {
				filePath += string(os.PathSeparator) + file.Name()
			}

			if GetConfig().File.CheckFormat(ext) && GetConfig().File.CheckOnly(file.Name()) {

				totalFile++

				// Execute with different threads
				ch := make(chan byte, 1)
				go func(f, p string) {
					line := Reader{f}.Search(p)
					yellow := color.New(color.FgYellow).SprintFunc()

					fmt.Print(yellow(f + ":" + strconv.Itoa(line) + " : "))

					if line > -1 {
						color.Green("True")
					} else {
						color.Red("False")
					}
					ch <- 1
				}(filePath, pattern)

				<-ch

			} else if file.IsDir() && GetConfig().File.Recursion {
				recFunc(filePath)
			}
		}
	}

	recFunc(directory)

	color.Yellow("\nTotal file: %v", totalFile)
}

// DIRECTORY
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
