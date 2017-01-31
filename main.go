package main

import (
	"bufio"
	"os"
	"log"
	"strings"
	//"github.com/jmoiron/sqlx"
)

type urlElement struct {
	url string
	reg *string
}

func main() {
	file, err := os.Open("rkn")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	urls := make([]urlElement, 0)
	for scanner.Scan() {
		tempText := scanner.Text()
		tempText = strings.Replace(tempText, "^", "", -1)
		tempText = strings.Replace(tempText, "$", "", -1)
		splittedText := strings.Split(tempText, "(")
		if len(splittedText) == 1 {
			urls = append(urls, urlElement{url:splittedText[0], reg:nil})
		} else {
			urls = append(urls, urlElement{url:splittedText[0], reg:&splittedText[1]})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
