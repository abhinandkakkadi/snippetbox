package main

import (
	"net/http"

	"github.com/justinas/alice"
)



func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Create a middleware chain containing our 'standard' middleware
  // which will be used for every request our application receives

	// Return the 'standard' middleware chain followed by the servemux
	standard := alice.New(app.recoverPanic,app.logRequest,secureHeaders)
	
	// secureHeaders middleware wraps the mux - every route registered to the router (mux) will get this middleware
	return standard.Then(mux)

}
