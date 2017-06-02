package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type server struct {
	Db *sqlx.DB
}

type UrlElement struct {
	Url string `db:"url"`
	Reg *string `db:"reg"`
}

type Url string

type UrlListTemplateData struct {
	UrlList []UrlElement
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

func (s *server) parseConfig(path string) {
	file, err := os.Open("rkn")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	urls := make([]UrlElement, 0)
	for scanner.Scan() {
		tempText := scanner.Text()
		tempText = strings.Replace(tempText, "^", "", -1)
		tempText = strings.Replace(tempText, "$", "", -1)
		splittedText := strings.Split(tempText, "(")
		if len(splittedText) == 1 {
			urls = append(urls, UrlElement{Url: splittedText[0], Reg: nil})
		} else {
			urls = append(urls, UrlElement{Url: splittedText[0], Reg: &splittedText[1]})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, singleUrl := range urls {
		fmt.Println(singleUrl.Url)
		if singleUrl.Reg != nil {
			*singleUrl.Reg = "(" + *singleUrl.Reg
			_, err = s.Db.NamedQuery("INSERT INTO `urls` (`url`, `reg`) VALUES (:url, :reg)", map[string]interface{}{"url": singleUrl.Url, "reg": *singleUrl.Reg, })
		} else {
			_, err = s.Db.NamedQuery("INSERT INTO `urls` (`url`, `reg`) VALUES (:url, :reg)", map[string]interface{}{"url": singleUrl.Url, "reg": singleUrl.Reg, })
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *server) addUrlToDbHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	tx := s.Db.MustBegin()
	for i := 1; i <= len(r.Form); i++ {
		tx.MustExec("INSERT INTO `urls` (`url`, `reg`) VALUES (?, ?)", r.PostFormValue("url"+strconv.Itoa(i)), "(/.*?)")
	}
	tx.Commit()
	http.Redirect(w, r, "/urlList/", 302)
}

func (s *server) updateUrlHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	tx := s.Db.MustBegin()
	for i := 1; i <= len(r.Form); i++ {
		tx.MustExec("UPDATE `urls` SET `url` = ? WHERE `url` = ?", r.PostFormValue("url"), r.PostFormValue("id"))
	}
	tx.Commit()
	http.Redirect(w, r, "/urlList/", 302)
}

func (s *server) deleteUrlHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	tx := s.Db.MustBegin()
	for i := 1; i <= len(r.Form); i++ {
		tx.MustExec("DELETE FROM `urls` WHERE `url` = ?", r.PostFormValue("url"))
	}
	tx.Commit()
	http.Redirect(w, r, "/urlList/", 302)
}

func (s *server) urlListHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	latexTemplate, err := template.ParseFiles("templates/urlList.tmpl.html")
	if err != nil {
		log.Println(err)
		return
	}
	tx, err := s.Db.Beginx()
	if err != nil {
		log.Println(err)
		return
	}
	urlList := make([]UrlElement, 0)
	if err := tx.Select(&urlList, "SELECT url, reg FROM urls ORDER BY id DESC"); err != nil {
		log.Println(err)
		return
	}
	err = latexTemplate.Execute(w, UrlListTemplateData{UrlList: urlList})
	if err != nil {
		log.Println(err)
		return
	}
	tx.Commit()
}

func (s *server) generateConfig() {
	_, err := os.Open("rkn_test_conf")
	if err != nil {
		log.Println(err)
	}
	data := make([]UrlElement, 0, 1000)
	s.Db.Select(&data, "SELECT DISTINCT url, reg FROM urls ORDER BY id DESC")
	log.Printf("%#v", data)
	for _, url := range data {
		if url.Reg != nil {
			fmt.Printf("^%s%s$\n", url.Url, *url.Reg)
		} else {
			fmt.Printf("^%s(/.*?)$\n", url.Url)
		}
	}
}

func main() {
	flag.Parse()
	loadConfig(*configFile)

	s := server{
		Db: sqlx.MustConnect("mysql", config.MysqlLogin+":"+config.MysqlPassword+"@tcp("+config.MysqlHost+")/"+config.MysqlDb+"?charset=utf8"),
	}
	s.generateConfig()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/addUrlToDb/", s.addUrlToDbHandler)
	http.HandleFunc("/deleteUrl/", s.deleteUrlHandler)
	http.HandleFunc("/updateUrl/", s.updateUrlHandler)
	http.HandleFunc("/urlList/", s.urlListHandler)

	log.Print("Server started at port 4002")
	err := http.ListenAndServe(":4002", nil)
	if err != nil {
		log.Fatal(err)
	}
}
