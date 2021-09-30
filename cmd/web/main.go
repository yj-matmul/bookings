package main

import (
	"encoding/gob"
	"flag"
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
var logFile *os.File
var dbInfoPath string

// main is the main application function
func main() {
	// dbInfoPath = "./static/db_info.txt"
	// dsn := loadDsn(dbInfoPath)

	infoLog, logFile = config.CustomLogger()
	app.InfoLog = infoLog
	defer logFile.Close()

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// parsing flags
	inProduction := flag.Bool("production", true, "application is in production")
	useCache := flag.Bool("cache", true, "use tempalte cache")
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbName := flag.String("dbname", "", "database name")
	dbUser := flag.String("dbuser", "", "database user")
	dbPassword := flag.String("dbpassword", "", "database password")
	dbPort := flag.String("dbport", "5001", "database port")
	dbSSL := flag.String("dbssl", "disable", "database ssl settings (disable, prefer, require")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		log.Println("Missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour // session의 유지 시간
	session.Cookie.Persist = true     // user가 browser를 종료해도 session을 유지한다는 option
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // true uses https, false uses http

	app.Session = session

	// connect to database
	app.InfoLog.Println("connect to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPassword, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	app.InfoLog.Println("connect to database!")

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

// loadDsn loads DB info from hidden file (not used)
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
