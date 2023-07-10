package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)





func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w,r)
		return
	}

	w.Write([]byte("Hello from abhinand"))

}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id,err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w,r)
		return
	}
	
	fmt.Fprintf(w,"snippet of ID %d",id)

}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create a new snippet..."))
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/",home)
	mux.HandleFunc("/snippet/view?id=2",snippetView)
	mux.HandleFunc("/snippet/create",snippetCreate)

	log.Print("starting server at port :4000")
  err := http.ListenAndServe(":4000",mux)
	log.Fatal(err)

}