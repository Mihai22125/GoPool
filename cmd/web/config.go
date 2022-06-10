package main

import (
	"html/template"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/mihai22125/goPool/pkg/models/mysql"

	"log"
)

type contextKey string

var contextKeyUser = contextKey("user")

type Config struct {
	AddrHTTPS string
	AddrHTTP  string
	StaticDir string
	Dsn       string
	Secret    string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	config        *Config
	session       *sessions.Session
	pools         *mysql.PoolModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
	machines      *mysql.MachineModel
	sessions      *mysql.SessionModel
	votes         *mysql.VoteModel
	apiKeys       *mysql.ApiKeyModel
	location      *time.Location
}
