# squid-ban-urls

Application for ban users in squid proxy server

db scheme stored in ./db.sql

## conf.json example

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

## Flags

```
Usage of /home/tochk/prj//bin/squid-ban-urls:
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
