package main

import (
	"fmt"
	"net/http"
)

// This middleware adds some of the headers proposed by OWASP to every response
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This header is tricky. You can block necessary resourses and break something
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// This func is application method and secureHeaders() isn't because we want to acces the logger
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// When one of the connections panics, it unwidns the stack (for now just calls any deffered funcs)
		// So we need to mark the connection as closed in the header
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If user is not authenticated, redirect to signin page
		// And return so no otherr middlewares are ran
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
			return
		}

		// This header means that pages requiring authentication
		// Are not stored in the browser cachep
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
