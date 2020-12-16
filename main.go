package main

import (
	"crypto/tls"
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
	mux.HandleFunc("/api/v1/order/orders", handlers.GetOrders).Methods("GET")
	mux.HandleFunc("/api/v1/order/ordersmails", handlers.SendMails).Methods("POST")
	mux.HandleFunc("/api/v1/order/updateOrder", handlers.UpDateOrders).Methods("GET")
	mux.HandleFunc("/api/v1/order/setstatus", handlers.SetStatus).Methods("POST")

	mux.HandleFunc("/api/v1/user/login", handlers.Login).Methods("GET", "POST")

	log.Println("The server is lisening")

	//log.Fatal(http.ListenAndServe(configs.GetPort(), mux))
	//log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", mux))
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))

}
