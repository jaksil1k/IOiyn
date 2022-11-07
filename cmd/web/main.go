package main

import (
	"IOiyn.kz/internal/models"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	games    *models.GameModel
	db       *models.DBModel
	users    *models.UserModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:password"+
		"@/game?multiStatements=true&parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		games:    &models.GameModel{DB: db},
		db:       &models.DBModel{DB: db},
		users:    &models.UserModel{DB: db},
	}

	err = app.db.DropTables()
	if err != nil {
		errorLog.Fatal(err)
	}
	err = app.db.CreateTables()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer app.db.DropTables()

	err = app.users.CreateInitialUsers()
	if err != nil {
		errorLog.Fatal(err)
	}

	err = app.games.CreateInitialGames()
	if err != nil {
		errorLog.Fatal(err)
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
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
