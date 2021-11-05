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

	recursion("", result, &recursions)

	if recursions.Stop {
		return recursions.LineFound
	} else {
		return -1
	}
}

func recursion(k string, i *interface{}, r *Recursions) {

	if r.Stop {
		return
	}

	switch v := (*i).(type) {
	case []interface{}:
		print(k, "[", *r)
		r.Depth++
		for kk, vv := range v {
			recursion(strconv.Itoa(kk), &vv, r)
		}
		r.Depth--
		print("", "]", *r)
	case map[string]interface{}:
		if k != "" {
			print(k, "{", *r)
			r.Depth++
		}
		for kk, vv := range v {
			recursion(kk, &vv, r)
		}
		if k != "" {
			r.Depth--
			print("", "}", *r)
		}
	default:
		print(k, v, *r)
		r.LineFound++
		if r.Pattern.hasKey(k) && r.Pattern.hasValue(v) {
			r.Stop = true
		}
	}
}

// Read YAML/JSON file
func readFile(path string) *interface{} {

	var result interface{}

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

func contain(s string, test interface{}) bool {
	return s == fmt.Sprintf("%v", test)
}

func (sp *SearchPattern) hasKey(test interface{}) bool {
	if !strings.Contains(sp.pattern, ":") {
		return false
	}

	s := strings.TrimSpace(strings.Split(sp.pattern, ":")[0])

	return contain(s, test)
}

func (sp *SearchPattern) hasValue(test interface{}) bool {
	if !strings.Contains(sp.pattern, ":") {
		return false
	}

	s := strings.TrimSpace(strings.Split(sp.pattern, ":")[1])

	return contain(s, test)
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
		log.Fatalf("[Unmarshal]: file %v; %v", path, err)
	}

	return err
}
