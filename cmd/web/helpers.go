package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer - This is done so that even if there is some runtime error, the user can get those message
	buf := new(bytes.Buffer)

	// write the template to buffer, instead of straight to http.ResponseWriter. If there's
	// an error , call our serverError() helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// if no error go ahead and set header to the intended header
	w.WriteHeader(status)

	// write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)

}

func (app *application) newTemplateData(r *http.Request) *templateData {

	return &templateData{
		CurrentYear: time.Now().Year(),
		// it will return the value corresponding to the given key
		// and also delete it
		//  If key does not exists, empty string will be returned
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		// add authentication status to template data
		IsAuthenticated: app.isAuthenticated(r),
	}

}

// create a decoder post form helper method
func (app *application) decodePostForm(r *http.Request, dst any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// for all other errors, return them as normal
		return err
	}

	return nil

}

// helper to check if the current user is an authenticated one
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(),"authenticatedUserID")
}