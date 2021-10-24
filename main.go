package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	readJson()
}

func readJson() {
	jsonFile, err := os.Open("C:\\DVL\\configTable\\src\\main\\resources\\tables\\regions\\3_tableMapping.json")

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
		for _, vv := range v {
			recursion("", vv)
		}
	case map[string]interface{}:
		fmt.Print(" ")
		for kk, vv := range v {
			recursion(kk, vv)
		}
	default:
		fmt.Println("-", k, ":", v)
	}
}
