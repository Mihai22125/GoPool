package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/mihai22125/goPool/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	cfg := new(Config)
	flag.StringVar(&cfg.AddrHTTPS, "addrHTTPS", ":4000", "HTTPS network address")
	flag.StringVar(&cfg.AddrHTTP, "addrHTTP", ":8000", "HTTP network address")

	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.Dsn, "dsn", "root:pass@tcp(database:3306)/goPool?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.Secret, "secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@sf", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(cfg.Secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	loc, _ := time.LoadLocation("Europe/Bucharest")

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		config:        cfg,
		session:       session,
		pools:         &mysql.PoolModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		machines:      &mysql.MachineModel{DB: db},
		sessions:      &mysql.SessionModel{DB: db},
		apiKeys:       &mysql.ApiKeyModel{DB: db},
		votes:         &mysql.VoteModel{DB: db},
		templateCache: templateCache,
		location:      loc,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         cfg.AddrHTTPS,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	restSrv := &http.Server{
		Addr:         cfg.AddrHTTP,
		ErrorLog:     errorLog,
		Handler:      app.restRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		infoLog.Printf("Starting server on %s\n", cfg.AddrHTTPS)
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		infoLog.Println("https: ", err)
	}()

	infoLog.Printf("Starting server on %s\n", cfg.AddrHTTP)
	err = restSrv.ListenAndServe()
	infoLog.Println("http: ", err)
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
