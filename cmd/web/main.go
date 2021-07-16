package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/yj-matmul/bookings/internal/config"
	"github.com/yj-matmul/bookings/internal/handlers"
	"github.com/yj-matmul/bookings/internal/models"
	"github.com/yj-matmul/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	err := run()

	fmt.Println(fmt.Sprintf("Starting application on prot %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false
	app.UseCache = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour // session의 유지 시간
	session.Cookie.Persist = true     // user가 browser를 종료해도 session을 유지한다는 option
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // true uses https, false uses http

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc

	render.NewTemplates(&app)

	Repo := handlers.NewRepo(&app)
	handlers.NewHandlers(Repo)

	return nil
}
