package handlers

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"victorydash/configs"
	"victorydash/models"

	goshopify "github.com/bold-commerce/go-shopify"
)

//GetOrdersFromShopify get the last 50 orders activitis from Shopify
func GetOrdersFromShopify() goshopify.OrdersResource {
	url := configs.GetUrlShopOrders()
	orders := goshopify.OrdersResource{}
	err := getJSON(url, &orders)

	if err != nil {
		fmt.Println(err)
	}
	return orders
}

//UpDateOrders get Orders from Shopify and update in our database
func UpDateOrders() {
	url := configs.GetUrlShopOrders()
	orders := goshopify.OrdersResource{}
	err := getJSON(url, &orders)

	if err != nil {
		fmt.Println(err)
	}

	for _, order := range orders.Orders {
		_, err2 := GetOrderByIDShopifyAndnameShopify(order.ID, order.Name)
		if err2 != nil {
			orderID := SaveOrder(order)
			fmt.Println(order.LineItems)
			for _, product := range order.LineItems {
				fmt.Println(product)
				InsertProduct(orderID, product)
			}
		} else {
			continue
		}
	}
}

//GetOrderByIDShopifyAndnameShopify return a order from DB by ID_Shopify  name_Shopify
func GetOrderByIDShopifyAndnameShopify(IDShopify int64, nameShopify string) (goshopify.Order, error) {
	order := goshopify.Order{}
	err := error(nil)
	query := `SELECT o.id_shopify FROM orders o
				WHERE o.id_shopify  LIKE ? and name_shopify  like ?`

	row, _ := configs.Query(query, IDShopify, nameShopify)
	if row.Next() {
		row.Scan(&order.ID)
	} else {
		err = errors.New("that Oreder no exist")
	}
	return order, err
}

//GetOrderByIDShopify return a order from DB by ID_Shopify
func GetOrderByIDShopify(IDShopify int64) (goshopify.Order, error) {
	order := goshopify.Order{}
	err := error(nil)
	query := `SELECT o.id_shopify, o.name_shopify FROM orders o
				WHERE o.id_shopify  LIKE ? `

	row, _ := configs.Query(query, IDShopify)
	if row.Next() {
		row.Scan(&order.ID, &order.Name)
	} else {
		err = errors.New("that Oreder no exist")
	}
	return order, err
}

//SaveOrder save on DB an Order
func SaveOrder(order goshopify.Order) int64 {
	query := `INSERT orders SET id_shopify =?, name_shopify =?, subtotal_price=?`
	result, _ := configs.Exec(query, order.ID, order.Name, order.SubtotalPrice)
	ID, _ := result.LastInsertId()
	return ID
}

//GetOrders return last 10 order
func GetOrders() (goshopify.OrdersResource, error) {
	orderResourse := goshopify.OrdersResource{}
	err := error(nil)
	query := `SELECT name_shopify,  id_shopify,send_provider, subtotal_price 
					FROM orders o 
					order by o.name_shopify desc
					limit 10`

	row, err := configs.Query(query)
	for row.Next() {
		order := goshopify.Order{}
		row.Scan(&order.Name, &order.ID, &order.Confirmed, &order.SubtotalPrice)
		orderResourse.Orders = append(orderResourse.Orders, order)
	}

	return orderResourse, err
}

//SendMails from arr of id_sopify send to
func SendMails(w http.ResponseWriter, r *http.Request) {
	var results map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &results); err != nil {
		log.Fatal(err)
	}
	out := fmt.Sprintf("%v", results["mails"])
	out1 := strings.TrimLeft(strings.TrimRight(out, "]"), "[")
	out2 := strings.Split(out1, " ")
	fmt.Println(out2)

	//crear el csv
	pf, err := CreateCsvOrderByProvider(out2, "WERM")
	if err != nil {
		//log.Panicln(err)
		results["No"] = "las ordenes selecionas no tienen productos para el proveedor  "
	}
	fmt.Println(pf)
	//mandar el csv adjunto en un correo
	configs.SendMailForWermProvider(pf)
	models.SendData(w, results)
}

func CreateCsvOrderByProvider(arrIdSopify []string, provider string) (string, error) {
	query := `SELECT o.name_shopify,TRIM(LEADING 'WERM-' FROM po.sku ) sku, po.quantity FROM orders  o  
				JOIN product_order po on o.id =po.order_id 
				WHERE  po.vendor LIKE ?
					AND o.id_shopify IN (%s)`
	ids := "'" + strings.Join(arrIdSopify[:], "','") + "'"
	query = fmt.Sprintf(query, ids)
	fmt.Println(query)
	rows, err := configs.Query(query, provider)
	if err != nil {
		log.Fatalf("In the Query..%s", err)
	}
	result := convertRowsInStringMatrix(rows)

	if len(result) == 0 {
		return "error", errors.New("no se lleno el arreglo ")
	}
	nameFile := writeCsvProvider(result)
	fmt.Println(nameFile)
	return nameFile, error(nil)
}
func writeCsvProvider(data [][]string) string {
	template := make([][]string, 0, 1000)
	titles := []string{"Bestellnummer", "Lineitem sku", "Lineitem quantity"}
	template = append(template, titles)
	for _, arr := range data {
		template = append(template, arr)
	}
	fmt.Println(template)
	date := time.Now().Format("2006.01.02 15:04:05")
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	nameFile := path + "/files/supplierOrders " + date + ".csv"

	newcsvFile, err := os.Create(nameFile)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(newcsvFile)
	_ = csvwriter.WriteAll(template)

	csvwriter.Flush()
	return nameFile
}
func convertRowsInStringMatrix(rows *sql.Rows) [][]string {
	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)
	cols, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	length := len(cols)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)
		if err != nil {
			panic(err.Error())
		}

		result = append(result, container)
	}
	return result
}
func getJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
