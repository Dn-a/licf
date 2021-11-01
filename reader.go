package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Reader struct {
	filePath string
}

func (r Reader) Search(pattern string) int {
	var result = readFile(r.filePath)

	var recursions = Recursions{0, 0, false, SearchPattern{pattern}}

	recursion("", *result, &recursions)

	if recursions.Stop {
		return recursions.LineFound
	} else {
		return -1
	}
}

func recursion(k string, i interface{}, r *Recursions) {

	if r.Stop {
		return
	}

	r.LineFound++

	switch v := i.(type) {
	case []interface{}:
		print(k, "[", *r)
		r.Depth++
		for kk, vv := range v {
			recursion(strconv.Itoa(kk), vv, r)
		}
		r.Depth--
		print("", "]", *r)
	case map[string]interface{}:
		if k != "" {
			print(k, "{", *r)
			r.Depth++
		}
		for kk, vv := range v {
			recursion(kk, vv, r)
		}
		if k != "" {
			r.Depth--
			print("", "}", *r)
		}
	default:
		print(k, v, *r)
		if r.Pattern.contain(k) && r.Pattern.contain(v.(string)) {
			r.Stop = true
		}
	}
}

// Read YAML/JSON file
func readFile(path string) *map[string]interface{} {

	var result map[string]interface{}

	unmarshalFile(path, &result)

	return &result
}

type Recursions struct {
	Depth     int
	LineFound int
	Stop      bool
	Pattern   SearchPattern
}

type SearchPattern struct {
	pattern string
}

func (r *Recursions) getSpace() string {
	var sb strings.Builder
	for i := 0; i < r.Depth; i++ {
		sb.WriteString("  ")
	}
	return sb.String()
}

func (sp *SearchPattern) contain(test string) bool {
	values := strings.Split(sp.pattern, ":")
	for _, val := range values {
		if val == test {
			return true
		}
	}
	return false
}

func print(k string, v interface{}, r Recursions) {
	if !GetConfig().File.Print {
		return
	}
	if k == "" && v != nil {
		fmt.Println(r.getSpace(), v)
	} else if v == nil {
		fmt.Println(r.getSpace(), k, ":")
	} else {
		fmt.Println(r.getSpace(),
			//"-",
			k, ":", v)
	}
}

// Unmarshal a file on strctured object
func unmarshalFile(path string, interf interface{}) error {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return err
		//log.Fatal(err)
	}

	fileType := strings.Replace(filepath.Ext(path), ".", "", -1)

	if fileType == "yaml" || fileType == "yml" {
		err = yaml.Unmarshal(file, interf)
	} else if fileType == "json" {
		err = json.Unmarshal(file, interf)
	} else {
		log.Fatalf("[Unmarshal]: type '%v' not allowed", fileType)
	}

	if err != nil {
		log.Fatalf("[Unmarshal]: %v", err)
	}

	return err
}
