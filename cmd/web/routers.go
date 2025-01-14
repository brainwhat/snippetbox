package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

// That's a stdlib implementation
// func (app *application) routes() http.Handler {
// 	mux := http.NewServeMux()

// 	// Here we bind our FileServer to the /static folder
// 	fileServer := http.FileServer(http.Dir("./ui/static/"))
// 	// So here let's say we request for /static/css/main.css
// 	// There is no /static folder inside /ui/static
// 	// So we need to strip this part
// 	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

// 	mux.HandleFunc("/", app.home)
// 	mux.HandleFunc("/snippet/view", app.snippetView)
// 	mux.HandleFunc("/snippet/create", app.snippetCreate)

// 	// Alice is a simple middleware chainer
// 	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

// 	return standard.Then(mux)

// }
