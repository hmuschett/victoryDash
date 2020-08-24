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
func UpDateOrders(w http.ResponseWriter, r *http.Request) {
	url := configs.GetUrlShopOrders()
	orders := goshopify.OrdersResource{}
	err := getJSON(url, &orders)
	results := make(map[string]string)

	if err != nil {
		fmt.Println(err)
		results["No"] = "Cant not read from shopify"
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

	results["SI"] = "Las ordernes fueron actualizadas"

	models.SendData(w, results)
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

//GetOrdersWERM return last 10 order
func GetOrdersWERM() (goshopify.OrdersResource, error) {
	orderResourse := goshopify.OrdersResource{}
	err := error(nil)
	query := `SELECT name_shopify,  id_shopify, send_provider, subtotal_price, status
				FROM orders o	
				WHERE o.id_shopify in (SELECT  id_shopify FROM  (select * from orders) as oo 
											join product_order po 
												on oo.id =po.order_id
										where po.vendor LIKE "%WERM%" )
				order by o.name_shopify desc
				limit 10`

	row, err := configs.Query(query)
	for row.Next() {
		order := goshopify.Order{}
		row.Scan(&order.Name, &order.ID, &order.Confirmed, &order.SubtotalPrice, &order.BrowserIp)
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
	fmt.Println("el nombre del csv es: " + pf)
	if err != nil {
		//log.Panicln(err)
		results["No"] = "las ordenes selecionas no tienen productos para el proveedor  "
	} else {
		//mandar el csv adjunto en un correo
		go configs.SendMailForWermProvider(pf)
		go updateStatusOrders(out2, "WERM")
	}

	models.SendData(w, results)
}

//SetStatus set the status order by order ID
func SetStatus(w http.ResponseWriter, r *http.Request) {
	var results map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &results); err != nil {
		log.Fatal(err)
	}
	arrIDSopify := fmt.Sprintf("%v", results["id"])
	fmt.Println(arrIDSopify)
	status := results["status"]
	fmt.Println(status)
	if status == "sendProvider" {
		go updateStatusOrder(arrIDSopify, "WERM", "sendProvider")
	} else if status == "received" {
		go updateStatusOrder(arrIDSopify, "WERM", "received")
	} else if status == "sendClient" {
		go updateStatusOrder(arrIDSopify, "WERM", "sendClient")
	}
}

func updateStatusOrders(arrIDSopify []string, provider string) {
	query := ` 	UPDATE orders 
					SET orders.status = 'sendProvider',
						orders.send_provider = TRUE 
				WHERE orders.id_shopify in (
							SELECT  id_shopify FROM (select * from orders) as  o  
								JOIN product_order po on o.id =po.order_id 
								WHERE  po.vendor LIKE CONCAT('%', ?, '%' )
									AND o.id_shopify IN (%s)
								)`
	ids := "'" + strings.Join(arrIDSopify[:], "','") + "'"
	query = strings.ReplaceAll(query, "%s", ids)
	configs.Exec(query, provider)
}

func updateStatusOrder(arrIDSopify string, provider string, status string) {
	query := ` 	UPDATE orders 
					SET orders.status = ?,
						orders.send_provider = TRUE 
				WHERE orders.id_shopify in (
							SELECT  id_shopify FROM (select * from orders) as  o  
								JOIN product_order po on o.id = po.order_id 
								WHERE  po.vendor LIKE CONCAT('%', ?, '%' )
									AND o.id_shopify = ?
								)`

	configs.Exec(query, status, provider, arrIDSopify)
}
func CreateCsvOrderByProvider(arrIdSopify []string, provider string) (string, error) {
	query := `SELECT o.name_shopify,SUBSTRING_INDEX(po.sku, '-', -1) sku, po.quantity FROM orders  o  
				JOIN product_order po on o.id =po.order_id 
				WHERE  po.vendor LIKE CONCAT('%', ?, '%' ) 
					AND o.id_shopify IN (%s)`
	ids := "'" + strings.Join(arrIdSopify[:], "','") + "'"
	query = strings.ReplaceAll(query, "%s", ids)
	fmt.Println(query)
	rows, err := configs.Query(query, provider)
	if err != nil {
		log.Fatalf("In the Query..%s", err)
	}
	result := convertRowsInStringMatrix(rows)

	fmt.Println(result)
	if len(result) == 0 {
		return "error no se lleno el csv", errors.New("no se lleno el arreglo ")
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
	nameFile := path + "/files/Bestellnummer " + date + ".csv"

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
