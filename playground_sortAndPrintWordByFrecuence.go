package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"regexp"
	"log"
)

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
	
	var text = "Once upon a midnight dreary, while I pondered weak and weary,Over many a quaint and curious volume of forgotten lore,While I nodded, nearly napping, suddenly there came a tapping,As of some one gently rapping, rapping at my chamber door."+
				"`'Tis some visitor,' I muttered, `tapping at my chamber door -" +
				"Only this, and nothing more.'"

	sortAndPrintWordByFrecuence(text)
		
}