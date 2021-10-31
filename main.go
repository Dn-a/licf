package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

const logo = `
_     _       __ 
| |   (_) ___ / _|
| |   | |/ __| |_ 
| |___| | (__|  _|
|_____|_|\___|_|  


`

func main() {

	initConfig()

	path := flag.String("path", "", "Specifies path of JSON file.")

	flag.Parse()

	files, err := ioutil.ReadDir(getCurrentDirectory())
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			fmt.Println(file.Name(), file.IsDir())
		}
	}

	p := "test.json"

	if *path != "" {
		p = *path
	}

	readJson(p)
}

func readJson(path string) {

	var result map[string]interface{}

	unmarshalFile(path, &result, tJSON)

	fmt.Println()

	recursion("", result)
}

func recursion(k string, i interface{}) {
	switch v := i.(type) {
	case []interface{}:
		print(k, "[")
		ws.Depth++
		for kk, vv := range v {
			recursion(strconv.Itoa(kk), vv)
		}
		ws.Depth--
		print("", "]")
	case map[string]interface{}:
		if k != "" {
			print(k, "{")
			ws.Depth++
		}
		for kk, vv := range v {
			recursion(kk, vv)
		}
		if k != "" {
			ws.Depth--
			print("", "}")
		}
	default:
		print(k, v)
	}
}

// SPACE
var ws WhiteSpace

type WhiteSpace struct {
	Depth int
}

func (ws *WhiteSpace) GetString() string {
	var sb strings.Builder
	for i := 0; i < ws.Depth; i++ {
		sb.WriteString("  ")
	}
	return sb.String()
}

func print(k string, v interface{}) {
	if k == "" && v != nil {
		fmt.Println(ws.GetString(), v)
	} else if v == nil {
		fmt.Println(ws.GetString(), k, ":")
	} else {
		fmt.Println(ws.GetString(),
			//"-",
			k, ":", v)
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

// Main Configuration
type configFile struct {
	Logo confifLogo
}

type confifLogo struct {
	Show bool `yaml:"show"`
}

type fileType string

const (
	tYAML fileType = "yaml"
	tJSON fileType = "json"
)

func initConfig() {
	var config configFile
	unmarshalFile("applications.yaml", &config, tYAML)

	if config.Logo.Show {
		color.Yellow(logo)
	}
}

// Unmarshal a file on strctured object
func unmarshalFile(path string, interf interface{}, tp fileType) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	if tp == tYAML {
		err = yaml.Unmarshal(file, interf)
	} else if tp == tJSON {
		err = json.Unmarshal(file, interf)
	} else {
		log.Fatalf("[Unmarshal]: type '%v' not allowed", tp)
	}

	if err != nil {
		log.Fatalf("[Unmarshal]: %v", err)
	}
}
