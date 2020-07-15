package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" //driver don connection
	"github.com/joho/godotenv"
)

type databaseConenction struct {
	username string
	password string
	host     string
	port     int
	dbName   string
}
type ClientMail struct {
	mail     string
	server   string
	password string
}

var (
	db           *sql.DB
	dbConenction *databaseConenction
	m            *ClientMail
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = strings.TrimRight(path, "/test")
	err = godotenv.Load(path + "/configs/.env")

	if err != nil {
		log.Println(err)
		log.Fatalf("Error loading .env file")
	}
	dbConenction = &databaseConenction{}
	dbConenction.username = os.Getenv("DB_USERNAME")
	dbConenction.password = os.Getenv("DB_PASSWORD")
	dbConenction.host = os.Getenv("DB_HOST")
	dbConenction.port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	dbConenction.dbName = os.Getenv("DB_DBNAME")

	m = &ClientMail{}
	m.mail = os.Getenv("CLIENT_MAIL")
	m.server = os.Getenv("SERVER_MAIL")
	m.password = os.Getenv("PASS_MAIL")

}

//CreateConnection to the Data Base
func CreateConnection() {
	if connetcion, err := sql.Open("mysql", generateConectionURL()); err != nil {
		panic(err)
	} else {
		db = connetcion
		log.Println("conecction db succeful!!")
	}
}

//CloseConnection for close de conection to db
func CloseConnection() {
	db.Close()
	log.Println("Close conecction db succeful!!")
}

//Ping make a ping to db
func Ping() {
	if err := db.Ping(); err != nil {
		panic((err))
	}
}

//GetPort return a web port from env
func GetPort() string {
	return os.Getenv("WEB_PORT")
}

//GetUrlShopOrders return the url from env
func GetUrlShopOrders() string {
	return os.Getenv("SHOP_ORDERS")
}

//GetMailConfig return
func GetMailConfig() ClientMail {
	return *m
}
func generateConectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConenction.username, dbConenction.password, dbConenction.host, dbConenction.port, dbConenction.dbName)
}

//Exec is the wrapper for db.exec to log is an error
func Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Println(err)
	}
	return result, err
}

//Query is the wrapper for db.Query to log is an error
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println(err)
	}
	return rows, err
}
