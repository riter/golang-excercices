package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"regexp"
	"log"
)

func getContentFile(fileName string) string {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func clearSpecialCharacters(text string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
        	log.Fatal(err)
    	}
	return reg.ReplaceAllString(text, " ")
}

func getFrecuenceByWord(words []string) map[string]int {
	output := make(map[string]int)
	for _, word := range words {
            output[word] = output[word] + 1
    	}
	return output
}

func printSortedByKey(wordsFrecuence map[string]int) {
	var keys []string
	for k := range wordsFrecuence {
	    keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
	    fmt.Println(k, wordsFrecuence[k])
	}
}

func sortAndPrintWordByFrecuence(text string) {
	text = clearSpecialCharacters(text)
	text = strings.ToLower(text)
	words := strings.Fields(text)
	
	output := getFrecuenceByWord(words)
	
	printSortedByKey(output)
}


func main() {
	
	var text = getContentFile("words.txt")

	sortAndPrintWordByFrecuence(text)
		
}