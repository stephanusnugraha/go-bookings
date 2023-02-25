package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/stephanusnugraha/go-bookings/pkg/config"
	"github.com/stephanusnugraha/go-bookings/pkg/handlers"
	"github.com/stephanusnugraha/go-bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8084"

var app config.AppConfig
var session *scs.SessionManager

// main app
func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	cache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = cache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
