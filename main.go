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

	//iterateFiles(*path, *pattern)

	folderIterator(*path, *pattern)

}

// Iterate all files from directory
func iterateFiles(directory, pattern string) {

	var totalFile int

	var recFunc func(dir string)

	recFunc = func(dir string) {

		if dir == "" {
			dir = getCurrentDirectory()
		}

		// Ignore direcotry or Files defined in config.yaml
		if lstF := strings.Split(dir, string(os.PathSeparator)); len(lstF) > 0 && GetConfig().File.CheckIgnore(lstF[len(lstF)-1]) {
			return
		}

		files, err := ioutil.ReadDir(dir)

		// if it is not a directory, checks if it is a file and append it to the files array,
		// otherwise it returns an error
		var isFile os.FileInfo
		var err2 error
		if err != nil {
			isFile, err2 = os.Stat(dir)
			files = append(files, isFile)
			if err2 != nil || !isFile.Mode().IsRegular() {
				log.Fatal(err)
				log.Fatal(err2)
			}
		}

		for _, file := range files {

			//time.Sleep(time.Microsecond)

			// Ignore direcotry or Files defined in config.yaml
			if GetConfig().File.CheckIgnore(file.Name()) {
				continue
			}

			ext := strings.Replace(filepath.Ext(file.Name()), ".", "", -1)

			fullPath := dir

			if isFile == nil || isFile.IsDir() {
				fullPath += string(os.PathSeparator) + file.Name()
			}

			// Check Format and allowed files, which are defined in config.yaml
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
				}(fullPath, pattern)

				<-ch

			} else if file.IsDir() && GetConfig().File.Recursion {
				recFunc(fullPath)
			}
		}
	}

	recFunc(directory)

	color.Yellow("\nTotal file: %v", totalFile)
}

// Current Directory
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getFile(dir, pattern string, file *os.FileInfo) {

	name := (*file).Name()

	// Ignore direcotry or Files defined in config.yaml
	if GetConfig().File.CheckIgnore(name) {
		return
	}

	ext := strings.Replace(filepath.Ext(name), ".", "", -1)

	fullPath := dir + string(os.PathSeparator) + name

	// Check Format and allowed files, which are defined in config.yaml
	if GetConfig().File.CheckFormat(ext) && GetConfig().File.CheckOnly(name) {

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
		}(fullPath, pattern)

		<-ch

	}
}

func folderIterator(dir, pattern string) {

	if dir == "" {
		dir = getCurrentDirectory()
	}

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		isFile, err2 := os.Stat(dir)
		if isFile.Mode().IsRegular() {
			getFile(dir, pattern, &isFile)
			return
		} else {
			log.Fatal(err, " ", err2)
		}
	}

	for _, file := range files {

		// Ignore direcotry or Files defined in config.yaml
		if file.IsDir() {

			if GetConfig().File.CheckIgnore(file.Name()) {
				continue
			} else if GetConfig().File.Recursion {
				fullPath := dir + string(os.PathSeparator) + file.Name()
				folderIterator(fullPath, pattern)
			}

		} else {
			getFile(dir, pattern, &file)
		}
	}
}
