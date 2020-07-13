package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"victorydash/configs"
	"victorydash/handlers"

	goshopify "github.com/bold-commerce/go-shopify"
	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

type data struct {
	Orders goshopify.OrdersResource
	Update func()
}

func (d *data) UpdateO() {
	fmt.Println("success!!")
}
func main() {
	configs.CreateConnection()
	defer configs.CloseConnection()
	defer fmt.Println("se cerro del sistema")

	mux := mux.NewRouter()
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	mux.PathPrefix("/assets/").Handler(s)

	Orders, _ := handlers.GetOrders()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(Data.Orders)
		templates.ExecuteTemplate(w, "index.html", Orders)
		//handlers.UpDateOrders()

	}).Methods("GET")
	mux.HandleFunc("/api/v1/ordersmails", handlers.SendMails).Methods("POST")
	mux.HandleFunc("/updateOrder", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("we are here :(")
		handlers.UpDateOrders()
	}).Methods("GET")

	log.Println("The server is lisening ")
	log.Fatal(http.ListenAndServe(configs.GetPort(), mux))

}
