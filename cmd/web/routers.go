package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// We change default NotFound func in httprouter to our global error
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// let's say we request for /static/css/main.css
	// There is no /static folder inside /ui/static
	// So we need to strip this part
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	// There is one problem with this. If i remove StripPrefix, there would be no error logged
	// Only browser can tell that there is no files on this path

	// We create dymanic middleware which doesn't involve static files
	dymanic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dymanic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dymanic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dymanic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dymanic.ThenFunc(app.snippetCreatePost))

	router.Handler(http.MethodGet, "/user/signup", dymanic.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dymanic.ThenFunc(app.userSignUpPost))
	router.Handler(http.MethodGet, "/user/signin", dymanic.ThenFunc(app.userSignIn))
	router.Handler(http.MethodPost, "/user/signin", dymanic.ThenFunc(app.userSignInPost))
	router.Handler(http.MethodPost, "/user/logout", dymanic.ThenFunc(app.userLogOutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
