package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// add security headers to every request
func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)

	})
}

// log request for every handler
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// recover panic middleware
func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// This deferred function will always be run in the event of a panic as Go unwinds a stack
		defer func() {

			// recover function check if there has been a panic or not
			// the err will of type panic whose underlying type is a string (passed to panic) or an error
			if err := recover(); err != nil {

				// Setting this Header trigger's the server to automatically close the connection
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))

			}
		}()

		next.ServeHTTP(w, r)

	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if user not authenticated, redirect them to login page
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// set "Cache-Control:no-store" header so that pages require authentication
		// are not stored in the users browser cache
		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// this function uses a customized CSRF cookie with the secure, path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler

}


func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		
		id := app.sessionManager.GetInt(r.Context(),"authenticatedUserID")
		
		// if authenticatedUserID does not exists
		if id == 0 {
			next.ServeHTTP(w,r)
			return
		}

		// Check if user with that ID exists in database
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w,err)
			return
		}

		// if ID exists, create a new copy of request (with isAuthenticatedContextKey value true in request context)
		if exists {
			ctx := context.WithValue(r.Context(),isAuthenticatedContextKey,true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w,r)

	})
}