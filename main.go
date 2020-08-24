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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Orders, _ := handlers.GetOrdersWERM()
		log.Println("Reload orders ")

		templates.ExecuteTemplate(w, "index.html", Orders)
	}).Methods("GET")
	mux.HandleFunc("/api/v1/order/ordersmails", handlers.SendMails).Methods("POST")
	mux.HandleFunc("/api/v1/order/updateOrder", handlers.UpDateOrders).Methods("GET")
	mux.HandleFunc("/api/v1/order/setstatus", handlers.SetStatus).Methods("POST")

	log.Println("The server is lisening")
	log.Fatal(http.ListenAndServe(configs.GetPort(), mux))

}
