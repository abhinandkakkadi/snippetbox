package main

import (
	"fmt"
	"net/http"
)

// add security headers to every request
func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		w.Header().Set("Content-Security-Policy","default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy","origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options","deny")
		w.Header().Set("X-XSS-Protection","0")
		
		next.ServeHTTP(w,r)

	})
}

// log request for every handler
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s",r.RemoteAddr,r.Proto,r.Method,r.URL.RequestURI())
		next.ServeHTTP(w,r)
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
				w.Header().Set("Connection","close")
				app.serverError(w,fmt.Errorf("%s",err))

			}
		}()

		next.ServeHTTP(w,r)

	})
}
