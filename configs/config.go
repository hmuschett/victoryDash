package configs

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	_ "github.com/denisenkom/go-mssqldb" //driver to MS SQLSERVER connection
	_ "github.com/go-sql-driver/mysql"   //driver to mySQL connection
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
	dbVic         *sql.DB
	dbSage        *sql.DB
	dbConenction  *databaseConenction
	dbSageConnect *databaseConenction
	m             *ClientMail
	appShop       goshopify.App
	clientShop    *goshopify.Client
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

	dbSageConnect = &databaseConenction{}
	dbSageConnect.username = os.Getenv("DB_SAGE_USERNAME")
	dbSageConnect.password = os.Getenv("DB_SAGE_PASSWORD")
	dbSageConnect.host = os.Getenv("DB_SAGE_HOST")
	dbSageConnect.port, _ = strconv.Atoi(os.Getenv("DB_SAGE_PORT"))
	dbSageConnect.dbName = os.Getenv("DB_SAGE_DBNAME")

	m = &ClientMail{}
	m.mail = os.Getenv("CLIENT_MAIL")
	m.server = os.Getenv("SERVER_MAIL")
	m.password = os.Getenv("PASS_MAIL")

	appShop = goshopify.App{
		ApiKey:   os.Getenv("SHOP_APIKEY"),
		Password: os.Getenv("SHOP_API_PASSWORD"),
	}
	clientShop = goshopify.NewClient(appShop, "victoryswitzerland", "", goshopify.WithVersion("2020-10"), goshopify.WithRetry(3))
}

//CreateConnection to the Vic Data Base
func CreateVicConnection() {
	if connetcion, err := sql.Open("mysql", generateConectionURLForVic()); err != nil {
		panic(err)
	} else {
		dbVic = connetcion
		log.Println("created connection with Vic DB succefuly!!")
	}
}

//CloseConnection for close de conection to Vic db
func CloseVicConnection() {
	dbVic.Close()
	log.Println("Close conecction dbVic succeful!!")
}

//Ping make a ping to Vic db
func PingToVic() {
	if err := dbVic.Ping(); err != nil {
		panic((err))
	} else {
		log.Println("conecction dbVic succeful!!")
	}
}

//CreateConnection to the Sage Data Base
func CreateSageConnection() {

	if connetcion, err := sql.Open("sqlserver", generateConectionURLForSage()); err != nil {
		log.Panicln(err)
	} else {
		dbSage = connetcion
		log.Println("created connection with SAGE DB succefuly!!")
	}
}

//CloseConnection for close de conection to Sage db
func CloseSageConnection() {
	if err := dbSage.Close(); err != nil {
		log.Println(err)
	}
	log.Println("Close conecction db SAGE succeful!!")
}

//Ping make a ping to Sage db
func PingToSage() {
	if err := dbSage.Ping(); err != nil {
		log.Println(err)
		//panic((err))
	} else {
		log.Println("conecction DB Sage succeful!!")
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

//GetCorsOrgin allowed from env
func GetCorsOrgin() string {
	return os.Getenv("CORS_ORIGIN")
}

//GetEnv envairoment from env
func GetEnv() string {
	return os.Getenv("ENV")
}

//GetRSAPath envairoment from env
func GetRSAPath() string {
	return os.Getenv("KEY_RSA_PATH")
}

//GetPathInAS2 envairoment from env
func GetPathInAS2() string {

	return os.Getenv("PATH_IN_AS2")
}

//GetServerAS2 envairoment from env
func GetServerAS2() string {
	return os.Getenv("SERVER_AS2")
}

//GetUserServerAS2  envairoment from env
func GetUserServerAS2() string {
	return os.Getenv("USER_SERVER_AS2")
}

//GetMailConfig return
func GetMailConfig() ClientMail {
	return *m
}

//GetClientShop sopify cliente connetion
func GetClientShop() *goshopify.Client {
	return clientShop
}

func generateConectionURLForVic() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbConenction.username, dbConenction.password, dbConenction.host, dbConenction.port, dbConenction.dbName)
}

func generateConectionURLForSage() string {
	query := url.Values{}
	query.Add("database", os.Getenv("DB_SAGE_DBNAME"))

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(dbSageConnect.username, dbSageConnect.password),
		Host:   fmt.Sprintf("%s:%d", dbSageConnect.host, dbSageConnect.port),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	return u.String()
}

//Exec is the wrapper for db.exec for Vic DB
func VicExec(query string, args ...interface{}) (sql.Result, error) {
	result, err := dbVic.Exec(query, args...)
	if err != nil {
		log.Println(err)
	}
	return result, err
}

//Query is the wrapper for db.Query for Vic DB
func VicQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := dbVic.Query(query, args...)
	if err != nil {
		log.Println(err)
	}
	return rows, err
}

//Exec is the wrapper for db.exec for SAGE DB
func SageExec(query string, args ...interface{}) (sql.Result, error) {
	result, err := dbSage.Exec(query, args...)
	if err != nil {
		log.Println(err)
	}
	return result, err
}

//Query is the wrapper for db.Query for SAGE DB
func SageQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := dbSage.Query(query, args...)
	if err != nil {
		log.Println(err)
	}
	return rows, err
}
