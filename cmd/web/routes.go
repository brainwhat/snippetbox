package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"snippetbox.brainwhat/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// We change default NotFound func in httprouter to our global error
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// We create dymanic middleware which doesn't involve static files
	dymanic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Unprotected routes
	router.Handler(http.MethodGet, "/", dymanic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/about", dymanic.ThenFunc(app.about))
	router.Handler(http.MethodGet, "/snippet/view/:id", dymanic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dymanic.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dymanic.ThenFunc(app.userSignUpPost))
	router.Handler(http.MethodGet, "/user/signin", dymanic.ThenFunc(app.userSignIn))
	router.Handler(http.MethodPost, "/user/signin", dymanic.ThenFunc(app.userSignInPost))

	protected := dymanic.Append(app.requireAuthentication)
	// Protected routes
	router.Handler(http.MethodGet, "/user/account", protected.ThenFunc(app.accountView))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogOutPost))
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
