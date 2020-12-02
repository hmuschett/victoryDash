package models

import (
	"errors"
	"victorydash/configs"
)

//User struct represent an user
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Users []User

var users = make(map[int]User)

//CreateUser in DB
func (user *User) CreateUser() *User {
	query := "INSERT users SET username=?, password=?"
	result, _ := configs.Exec(query, user.Username, user.Password)
	user.ID, _ = result.LastInsertId()
	return user
}

//GetUserByID FROM vic
func GetUserByID(ID int64) (User, error) {
	user := User{}
	err := error(nil)
	query := "SELECT id, username, password FROM users u WHERE u.id=?"
	rows, _ := configs.Query(query, ID)
	if rows.Next() {
		rows.Scan(&user.ID, &user.Username, &user.Password)
	} else {
		err = errors.New("El usuario no existe")
	}
	return user, err
}

//GetUserByUsername
func GetUserByUsername(username string) (User, error) {
	user := User{}
	err := error(nil)
	query := "SELECT id, username, password FROM users u WHERE u.username=?"
	rows, _ := configs.Query(query, username)
	if rows.Next() {
		rows.Scan(&user.ID, &user.Username, &user.Password)
	} else {
		err = errors.New("El usuario no existe")
	}
	return user, err
}

//Login in to app
func Login(username, password string) (User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return User{}, errors.New("User or Password unvalidate")
	}
	return user, nil
}

//GetUsers get all users FROM vic
func GetUsers() Users {
	users := Users{}
	query := "SELECT id, username, password FROM users"
	rows, _ := configs.Query(query)
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Username, &user.Password)
		users = append(users, user)
	}
	return users
}

//UpDateUser in DB
func UpDateUser(user User) error {
	query := "UPDATE users u SET u.username=?, u.password=? where u.id =?"
	_, err := configs.Query(query, user.Username, user.Password, user.ID)
	return err
}

//DeleteUser in DB
func DeleteUser(ID int) error {
	query := "DELETE FROM users u where u.id =?"
	_, err := configs.Query(query, ID)
	return err
}
