package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"victorydash/models"
	"victorydash/utils"

	uuid "github.com/satori/go.uuid"
)

//Login function login in to the app
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookie("victory_session"))
	utils.EnableCors(w)

	if r.Method == "POST" {
		user := &models.User{}
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &user); err != nil {
			log.Fatal(err)
		}

		fmt.Println(user)
		_, err2 := models.Login(user.Username, user.Password)
		if err2 != nil {
			fmt.Println(err2)
		} else {
			//utils.SetSession(&userRes, w)
			cookieValue := utils.Cookie{
				Name: "victory_session", Value: uuid.NewV4().String(),
			}

			models.SendData(w, cookieValue)
			fmt.Println("autenticate succefull")
		}
	}
}
