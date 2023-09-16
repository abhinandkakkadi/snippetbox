package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	// custom handler for 404 Not Found response
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Now there will be a response body for 404 response status
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives

	// Return the 'standard' middleware chain followed by the servemux
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// secureHeaders middleware wraps the mux - every route registered to the router (mux) will get this middleware
	return standard.Then(router)

}
