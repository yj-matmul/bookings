package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/yj-matmul/bookings/internal/config"
	"github.com/yj-matmul/bookings/internal/driver"
	"github.com/yj-matmul/bookings/internal/handlers"
	"github.com/yj-matmul/bookings/internal/helpers"
	"github.com/yj-matmul/bookings/internal/models"
	"github.com/yj-matmul/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var dbInfoPath string

// main is the main application function
func main() {
	dbInfoPath = "./static/db_info.txt"
	dsn := loadDsn(dbInfoPath)
	db, err := run(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Starting application on prot %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run(dsn string) (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false
	app.UseCache = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour // session의 유지 시간
	session.Cookie.Persist = true     // user가 browser를 종료해도 session을 유지한다는 option
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // true uses https, false uses http

	app.Session = session

	// connect to database
	log.Println("connect to database...")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("connect to database!")

	// create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

func loadDsn(path string) string {
	// loading password of DB
	var password string

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open db info file")
	}
	defer file.Close()
	fmt.Fscan(file, &password)
	log.Println("Loaded password from db info file")

	return fmt.Sprintf("host=localhost port=5001 dbname=bookings user=postgres password=%s", password)
}
