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
	"golang.org/x/crypto/acme/autocert"
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
	configs.CreateVicConnection()
	defer configs.CloseVicConnection()

	configs.CreateSageConnection()
	defer configs.CloseSageConnection()

	mux := mux.NewRouter()
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	mux.PathPrefix("/assets/").Handler(s)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Orders, _ := handlers.GetOrdersWERM()
		log.Println("Reload orders ")

		templates.ExecuteTemplate(w, "index.html", Orders)
	}).Methods("GET")

	mux.HandleFunc("/api/v1/user/login", handlers.Login).Methods("POST")

	mux.HandleFunc("/api/v1/order/orders", handlers.GetOrders).Methods("GET")
	mux.HandleFunc("/api/v1/order/ordersmails", handlers.SendMails).Methods("POST")
	mux.HandleFunc("/api/v1/order/updateOrder", handlers.UpDateOrders).Methods("GET")
	mux.HandleFunc("/api/v1/order/setstatus", handlers.SetStatus).Methods("POST")

	mux.HandleFunc("/api/v1/posorder/orders", handlers.GetPOSOrders).Methods("GET")
	mux.HandleFunc("/api/v1/posorder/ordersmails", handlers.SendMailPOSOrders).Methods("POST")
	mux.HandleFunc("/api/v1/posorder/refoundorder", handlers.SendMailPOSRefoundOrders).Methods("POST")

	if configs.GetEnv() == "dev" {
		log.Fatal(http.ListenAndServe(configs.GetPort(), mux))
	} else if false {

		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				/* 	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA, */

				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			},
			NextProtos: []string{"h2", "http/1.1"},
		}
		srv := &http.Server{
			Addr:      ":443",
			Handler:   mux,
			TLSConfig: cfg,
			//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
	} else {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("api.victoryswitzerland.com"), //Your domain here
			Cache:      autocert.DirCache("certs"),                           //Folder for storing certificates
		}
		server := &http.Server{
			Addr:    ":https",
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}
		// serve HTTP, which will redirect automatically to HTTPS
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

		log.Fatal(server.ListenAndServeTLS("", "")) //Key and cert are coming from Let's Encrypt
	}
	log.Println("The server is lisening")
}
