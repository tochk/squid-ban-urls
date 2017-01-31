package main

import (
	"bufio"
	"os"
	"log"
	"strings"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type urlElement struct {
	url string `db:"url"`
	reg *string `db:"reg"`
}

var (
	configFile = flag.String("Config", "conf.json", "Where to read the Config from")
)

var config struct {
	MysqlLogin    string `json:"mysqlLogin"`
	MysqlPassword string `json:"mysqlPassword"`
	MysqlHost     string `json:"mysqlHost"`
	MysqlDb       string `json:"mysqlDb"`
}

func loadConfig(path string) error {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, &config)
}

func main() {
	flag.Parse()
	loadConfig(*configFile)

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

	dataBase, err := sqlx.Connect("mysql", config.MysqlLogin+":"+config.MysqlPassword+"@tcp("+config.MysqlHost+")/"+config.MysqlDb+"?charset=utf8")
	defer dataBase.Close()
	if err != nil {
		log.Print(err)
	}
	for _, singleUrl := range urls {
		fmt.Println(singleUrl.url)
		if singleUrl.reg != nil {
			*singleUrl.reg = "(" + *singleUrl.reg
			_, err = dataBase.NamedQuery("INSERT INTO `urls` (`url`, `reg`) VALUES (:url, :reg)", map[string]interface{}{ "url": singleUrl.url, "reg": *singleUrl.reg,})
		} else {
			_, err = dataBase.NamedQuery("INSERT INTO `urls` (`url`, `reg`) VALUES (:url, :reg)", map[string]interface{}{ "url": singleUrl.url, "reg": singleUrl.reg,})
		}

		if err != nil {
			log.Fatal(err)
		}
	}

}
