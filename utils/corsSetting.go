package utils

import (
	"net/http"
	"victorydash/configs"
)

func EnableCors(w *http.ResponseWriter) {
	//(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Origin", configs.GetCorsOrgin())
	//(*w).Header().Set("Access-Control-Allow-Origin", "https://dashfronter-2-e6ygg.ondigitalocean.app")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
