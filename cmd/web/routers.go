package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// That's a stdlib implementation
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Here we bind our FileServer to the /static folder
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// So here let's say we request for /static/css/main.css
	// There is no /static folder inside /ui/static
	// So we need to strip this part
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Alice is a simple middleware chainer
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)

}
