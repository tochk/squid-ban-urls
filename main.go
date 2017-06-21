package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"gopkg.in/ldap.v2"
)

type server struct {
	Db *sqlx.DB
}

type UrlElement struct {
	Id  int `db:"id"`
	Url string `db:"url"`
	Reg *string `db:"reg"`
}

type Url string

var (
	configFile     = flag.String("Config", "conf.json", "Where to read the Config from")
	servicePort    = flag.Int("Port", 4002, "Application port")
	configFilePath = flag.String("ConfigFilePath", "rkn_test_conf", "Config file path")
	store          = sessions.NewCookieStore([]byte("applicationDataLP"))
)

var config struct {
	MysqlLogin    string `json:"mysqlLogin"`
	MysqlPassword string `json:"mysqlPassword"`
	MysqlHost     string `json:"mysqlHost"`
	MysqlDb       string `json:"mysqlDb"`
	LdapUser      string `json:"ldapUser"`
	LdapPassword  string `json:"ldapPassword"`
}

func loadConfig(path string) error {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, &config)
}

func (s *server) addUrlToDbHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	if session.Values["userName"] == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
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

func (s *server) deleteUrlHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	if session.Values["userName"] == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	urlId := r.URL.Path[len("/deleteUrl/"):]
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = s.Db.Exec("DELETE FROM `urls` WHERE `id` = ?", urlId)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/urlList/", 302)
}

func (s *server) urlListHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	if session.Values["userName"] == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	latexTemplate, err := template.ParseFiles("templates/urlList.tmpl.html")
	if err != nil {
		log.Println(err)
		return
	}
	urlList := make([]UrlElement, 0)
	if err := s.Db.Select(&urlList, "SELECT id, url, reg FROM urls ORDER BY id DESC"); err != nil {
		log.Println(err)
		return
	}
	err = latexTemplate.Execute(w, urlList)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *server) reload() {
	for {
		st := time.Now()
		acl, err := s.generateConfig()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 30)
			continue
		}
		log.Printf("Config build finished in %s", time.Now().Sub(st))
		oldconf, err := ioutil.ReadFile(*configFilePath)
		if bytes.Equal([]byte(acl), oldconf) {
			log.Println("Old config matches new")
			time.Sleep(time.Second * 30)
			continue
		}
		log.Println("Writing new config")
		file, err := os.Create(*configFilePath)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 30)
			continue
		}
		_, err = file.WriteString(acl)
		if err != nil {
			log.Println(err)
		}
		file.Close()
		time.Sleep(time.Second * 30)
	}
}

func (s *server) generateConfig() (acl string, err error) {
	data := make([]UrlElement, 0)
	err = s.Db.Select(&data, "SELECT DISTINCT url FROM urls ORDER BY id DESC")
	if err != nil {
		return
	}
	for _, url := range data {
		acl += fmt.Sprintf("^%s(.*?)$\n", strings.Replace(url.Url, "/", "\\/", -1))
	}
	log.Println("Config generated successfuly")
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	if session.Values["userName"] != nil {
		http.Redirect(w, r, "/add/", 302)
		return
	}
	latexTemplate, err := template.ParseFiles("templates/index.tmpl.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = latexTemplate.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *server) addHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	if session.Values["userName"] == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	latexTemplate, err := template.ParseFiles("templates/add.tmpl.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = latexTemplate.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func auth(login, password string) (username string, err error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "main.sgu.ru", 389))
	if err != nil {
		return
	}
	defer l.Close()

	err = l.Bind(config.LdapUser, config.LdapPassword)
	if err != nil {
		return
	}

	searchRequest := ldap.NewSearchRequest(
		"dc=main,dc=sgu,dc=ru",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(sAMAccountName="+login+"))",
		[]string{"cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return
	}

	if len(sr.Entries) == 1 {
		username = sr.Entries[0].GetAttributeValue("cn")
	} else {
		err = errors.New("User not found")
		return
	}

	err = l.Bind(username, password)
	if err != nil {
		return
	}

	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Loaded login page from " + r.RemoteAddr)
	r.ParseForm()
	session, _ := store.Get(r, "applicationData")
	userName, err := auth(r.Form["login"][0], r.Form["password"][0])
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		session, _ = store.Get(r, "applicationData")
		session.Values["userName"] = userName
		session.Save(r, w)
		http.Redirect(w, r, "/add/", 302)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Loaded logout page from " + r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	session, _ = store.Get(r, "applicationData")
	session.Values["userName"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func main() {
	flag.Parse()
	err := loadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded from " + *configFile)

	s := server{
		Db: sqlx.MustConnect("mysql", config.MysqlLogin+":"+config.MysqlPassword+"@tcp("+config.MysqlHost+")/"+config.MysqlDb+"?charset=utf8"),
	}
	defer s.Db.Close()
	log.Printf("Connected to database on %s", config.MysqlHost)

	go s.reload()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add/", s.addHandler)
	http.HandleFunc("/addUrlToDb/", s.addUrlToDbHandler)
	http.HandleFunc("/deleteUrl/", s.deleteUrlHandler)
	http.HandleFunc("/urlList/", s.urlListHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logoutHandler)
	log.Print("Server started at port " + strconv.Itoa(*servicePort))
	err = http.ListenAndServe(":"+strconv.Itoa(*servicePort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
