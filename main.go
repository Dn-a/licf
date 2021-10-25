package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	readJson()
}

func readJson() {
	//jsonFile, err := os.Open("C:\\DVL\\configTable\\src\\main\\resources\\tables\\regions\\3_tableMapping.json")
	jsonFile, err := os.Open("test.json")

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	recursion("", result)

}

func recursion(k string, i interface{}) {
	switch v := i.(type) {
	case []interface{}:
		fmt.Println(ws.String(), "-", k, ":")
		ws.Size++
		for _, vv := range v {
			recursion("", vv)
			fmt.Println()
		}
		ws.Size--
	case map[string]interface{}:
		if k != "" {
			fmt.Println(ws.String(), "-", k, ":")
			ws.Size++
		}
		for kk, vv := range v {
			recursion(kk, vv)
		}
		if k != "" {
			ws.Size--
		}
	default:
		fmt.Println(ws.String(), "-", k, ":", v)
	}
}

var ws WhiteSpace

type WhiteSpace struct {
	Size int
}

func (ws *WhiteSpace) String() string {
	var sb strings.Builder
	for i := 0; i < ws.Size; i++ {
		sb.WriteString("  ")
	}
	return sb.String()
}
