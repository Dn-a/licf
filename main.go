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
)

const logo = `
_     _       __ 
| |   (_) ___ / _|
| |   | |/ __| |_ 
| |___| | (__|  _|
|_____|_|\___|_|  


`

func main() {
	color.Yellow(logo)

	p := flag.String("p", "", "Specifies path of JSON file.")

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

	path := "test.json"

	if *p != "" {
		fmt.Println(p)
		path = *p
	}

	readJson(path)
}

func readJson(path string) {
	//jsonFile, err := os.Open("C:\\DVL\\configTable\\src\\main\\resources\\tables\\regions\\3_tableMapping.json")
	jsonFile, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}

	json.Unmarshal([]byte(byteValue), &result)

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

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
