package main

import (
	//"bufio"
	//"os"
	"log"
	//"strings"
	//"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"io/ioutil"
	"encoding/json"
	//"fmt"
	"net/http"
)

type urlElement struct {
	Url string `db:"url"`
	Reg *string `db:"reg"`
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

	/*file, err := os.Open("rkn")
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
			urls = append(urls, urlElement{Url: splittedText[0], Reg: nil})
		} else {
			urls = append(urls, urlElement{Url: splittedText[0], Reg: &splittedText[1]})
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
		fmt.Println(singleUrl.Url)
		if singleUrl.Reg != nil {
			*singleUrl.Reg = "(" + *singleUrl.Reg
			_, err = dataBase.NamedQuery("INSERT INTO `urls` (`url`, `reg`) VALUES (:url, :reg)", map[string]interface{}{ "url": singleUrl.Url, "reg": *singleUrl.Reg,})
		} else {
			_, err = dataBase.NamedQuery("INSERT INTO `urls` (`url`, `reg`) VALUES (:url, :reg)", map[string]interface{}{ "url": singleUrl.Url, "reg": singleUrl.Reg,})
		}

		if err != nil {
			log.Fatal(err)
		}
	}*/



	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Server started at port 4002")
	err := http.ListenAndServe(":4002", nil)
	if err != nil {
		log.Fatal(err)
	}
}
