package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const baseUrl = "ClinicalTrials.gov/api/"
const fmtUrl = "&fmt=json"
const fldUrl = "&fields="
const queryUrl = "query/study_fields?expr="

var url = baseUrl
var query string
var fields string

func main() {
	args := os.Args[1:]
	for idx, arg := range args {
		switch arg {
		case "-p", "--print-options":
			printOptions()
			os.Exit(0)
		case "-q", "--query":
			query = strings.Join(strings.Fields(args[idx+1]), "+")
			url += queryUrl + query
			fmt.Println(url)
		case "-o", "--options":
			fields = strings.Join(strings.Fields(args[idx+1]), ",")
			url += fldUrl + fields
			fmt.Println(url)
		}
	}
	data := retrieveJson(url, true)
	fmt.Println(data)
}

//map[string]interface{}
func retrieveXml(url string) string {
	var rawData string
	err := requests.
		URL(url).
		ToString(&rawData).
		Fetch(context.Background())
	if err != nil {
		panic("An error occurred in the request.")
	}
	//xml := strings.NewReader(rawData)
	//jsonBuffer, err := xml2json.Convert(xml)
	//if err != nil {
	//	panic("An error occurred in the parsing.")
	//}
	//jsonBytes, err := ioutil.ReadAll(jsonBuffer)
	//if err != nil {
	//	panic("An error occured while debuffering bytes.")
	//}
	//var data map[string]interface{}
	//if err := json.Unmarshal(jsonBytes, &data); err != nil {
	//	panic(err)
	//}
	return rawData
}

func retrieveJson(url string, save bool) interface{} {
	url += fmtUrl
	fmt.Println(url)
	var data interface{}
	err := requests.
		URL(url).
		ToJSON(&data).
		Fetch(context.Background())
	if err != nil {
		panic("An error occurred in the request.")
	}
	if save {
		file, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			panic("An error occurred while marshalling.")
		}
		err = ioutil.WriteFile("test.json", file, 0644)
		if err != nil {
			panic("An error occurred while writing the file.")
		}
	}
	return data
}

func printOptions() {
	url := "https://clinicaltrials.gov/api/info/study_fields_list"
	data := retrieveXml(url)
	re := regexp.MustCompile(`"[^"]+"`)
	options := re.FindAllString(data, -1)
	for idx, val := range options {
		fmt.Println(idx, val)
	}
}
