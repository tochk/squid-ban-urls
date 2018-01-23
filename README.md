# squid-ban-urls

Application for ban URL's in squid proxy server

db scheme stored in [db.sql](https://github.com/tochk/squid-ban-urls/blob/master/db.sql)

## Installing and configuring

For install application:

1) [Download](https://golang.org) and install go 1.9+

2) Install quick template

```
go get -u github.com/valyala/quicktemplate
go get -u github.com/valyala/quicktemplate/qtc
```

3) Download squid-ban-urls application
```
go get -u github.com/tochk/squid-ban-urls
```

4) Build application
```
cd $GOPATH/src/github.com/tochk/squid-ban-urls
$GOPATH/bin/qtc
go install github.com/tochk/squid-ban-urls
```

5) Make conf.json file and use flags you need

6) Run application

```
$GOPATH/bin/squid-ban-urls
```



### conf.json example

Don't forget about secure session key

```
{
  "mysqlLogin": "root",
  "mysqlPassword": "",
  "mysqlHost": "localhost:3306",
  "mysqlDb": "sample",
  "ldapUser": "admin",
  "ldapPassword": "password",
  "ldapBaseDN": "OU=Test,OU=Group,DC=example,DC=com",
  "ldapAddress": "example.com",
  "ldapPort": 389,
  "sessionKey": "GENERATE KEY FOR SESSIONS"
}
```

### Flags

```
Usage of squid-ban-urls:
  -config string
        Where to read the config from (default "conf.json")
  -per_page int
        URL's per page (default 50)
  -port int
        Application port (default 4002)
  -restart_interval int
        Squid restart interval (seconds) (default 30)
  -squid_config_path string
        Config file path (default "squid_acl")
```
