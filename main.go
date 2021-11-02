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

	files := Files{*path, *pattern, 0}
	files.iterateFiles()

}

type Files struct {
	path      string
	pattern   string
	totalFile int
}

// Iterate all files from directory
func (f *Files) iterateFiles() {

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
			if !f.Mode().IsRegular() || err2 != nil {
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

				f.totalFile++

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
				}(filePath, f.pattern)

				<-ch

			} else if file.IsDir() && GetConfig().File.Recursion {
				recFunc(filePath)
			}
		}
	}

	recFunc(f.path)

	color.Yellow("\nTotal file: %v", f.totalFile)
}

// DIRECTORY
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
