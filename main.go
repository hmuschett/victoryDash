package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	mux := mux.NewRouter()
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	mux.PathPrefix("/assets/").Handler(s)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintln(w, "hello")
		templates.ExecuteTemplate(w, "index.html", nil)
	}).Methods("GET")

	log.Println("The server is lisening on 3000 port")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	log.Println(port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

}
