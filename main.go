package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-systemd/dbus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/tochk/squid-ban-urls/templates"
	"gopkg.in/ldap.v2"
)

type server struct {
	Db *sqlx.DB
	dc *dbus.Conn
}

type UrlElement = templates.UrlElement

type Pagination = templates.Pagination

var (
	configFile     = flag.String("Config", "conf.json", "Where to read the Config from")
	servicePort    = flag.Int("Port", 4002, "Application port")
	configFilePath = flag.String("ConfigFilePath", "squid_acl", "Config file path")
	store          *sessions.CookieStore
)

var config struct {
	MysqlLogin    string `json:"mysqlLogin"`
	MysqlPassword string `json:"mysqlPassword"`
	MysqlHost     string `json:"mysqlHost"`
	MysqlDb       string `json:"mysqlDb"`
	LdapUser      string `json:"ldapUser"`
	LdapPassword  string `json:"ldapPassword"`
	LdapBaseDN    string `json:"ldapBaseDN"`
	SessionKey    string `json:"sessionKey"`
}

func loadConfig(path string) error {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return err
	}
	store = sessions.NewCookieStore([]byte(config.SessionKey))
	return nil
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
	url := strings.Split(r.URL.Path, "/")
	var urlList []UrlElement
	var pagination Pagination
	if url[2] == "page" {
		page, err := strconv.Atoi(url[3])
		if err != nil {
			log.Println(err)
			return
		}
		pagination = s.paginationCalc(page, 50)
	} else {
		pagination = s.paginationCalc(1, 50)
	}
	if err := s.Db.Select(&urlList, "SELECT id, url, reg FROM urls ORDER BY id DESC LIMIT ? OFFSET ?", 50, pagination.Offset); err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, templates.ListPage(pagination, urlList))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	switch r.URL.Path {
	case "/login/":
		r.ParseForm()
		userName, err := auth(r.PostForm.Get("login"), r.PostForm.Get("password"))
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", 302)
		} else {
			session, _ = store.Get(r, "applicationData")
			session.Values["userName"] = userName
			session.Save(r, w)
			http.Redirect(w, r, "/add/", 302)
		}
	case "/logout/":
		session.Values["userName"] = nil
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
	default:
		if session.Values["userName"] != nil {
			http.Redirect(w, r, "/add/", 302)
			return
		}
		fmt.Fprint(w, templates.LoginPage())
	}
}

func (s *server) addHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Loaded %s page from %s", r.URL.Path, r.RemoteAddr)
	session, _ := store.Get(r, "applicationData")
	if session.Values["userName"] == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	fmt.Fprint(w, templates.AddPage())
}

func auth(login, password string) (username string, err error) {
	if password == "" {
		return "", errors.New("empty password")
	}
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
		config.LdapBaseDN,
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
		err = errors.New("user not found")
		return
	}

	err = l.Bind(username, password)
	if err != nil {
		return
	}

	return
}

func (s *server) reload() error {
	st := time.Now()
	newconf, err := s.generateConfig()
	if err != nil {
		return err
	}
	log.Printf("Config build finished in %s", time.Now().Sub(st))
	oldconf, err := ioutil.ReadFile(*configFilePath)
	if bytes.Equal([]byte(newconf), oldconf) {
		log.Println("Old config matches new")
		return nil
	}
	log.Println("Writing new config")
	if err = ioutil.WriteFile(*configFilePath, []byte(newconf), os.ModePerm); err != nil {
		return err
	}
	log.Printf("Write finished in %s", time.Now().Sub(st))
	ch := make(chan string, 1)
	log.Println("Restarting squid")
	st = time.Now()
	id, err := s.dc.ReloadOrTryRestartUnit("squid.service", "fail", ch)
	if err != nil {
		return err
	}
	log.Printf("Job id: %d", id)
	res := <-ch
	log.Printf("Result: %s in %s", res, time.Now().Sub(st))
	return nil
}

func (s *server) run() {
	for {
		if err := s.reload(); err != nil {
			log.Printf("reload error: %s", err)
		}
		time.Sleep(time.Second * 30)
	}
}

func (s *server) generateConfig() (string, error) {
	var data []UrlElement
	if err := s.Db.Select(&data, "SELECT DISTINCT url FROM urls ORDER BY id DESC"); err != nil {
		return "", err
	}
	r := make([]string, 0, len(data))
	for _, v := range data {
		r = append(r, "^"+regexp.QuoteMeta(v.Url)+".*$")
	}
	return strings.Join(r, "\n"), nil
}

func (s *server) paginationCalc(page, perPage int) Pagination {
	var (
		count      int
		pagination Pagination
		err        error
	)
	if page < 1 {
		page = 1
	}
	pagination.CurrentPage = page
	pagination.PerPage = perPage
	pagination.Offset = perPage * (page - 1)
	err = s.Db.Get(&count, "SELECT COUNT(*) FROM urls")

	if err != nil {
		log.Println(err)
		return Pagination{}
	}
	if count > perPage*page {
		pagination.NextPage = pagination.CurrentPage + 1
		if pagination.NextPage != (count/perPage)+1 {
			pagination.LastPage = (count / perPage) + 1
		}
	}
	if pagination.CurrentPage > 1 {
		pagination.PrevPage = pagination.CurrentPage - 1
	}
	return pagination
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

	//if s.dc, err = dbus.New(); err != nil {
	//	log.Fatal(err)
	//}

	//go s.run()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add/", s.addHandler)
	http.HandleFunc("/addUrlToDb/", s.addUrlToDbHandler)
	http.HandleFunc("/deleteUrl/", s.deleteUrlHandler)
	http.HandleFunc("/urlList/", s.urlListHandler)
	log.Print("Server started at port " + strconv.Itoa(*servicePort))
	err = http.ListenAndServe(":"+strconv.Itoa(*servicePort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
