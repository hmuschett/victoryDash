package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"victorydash/configs"
	"victorydash/models"
	"victorydash/utils"

	goshopify "github.com/bold-commerce/go-shopify/v3"

	"github.com/shopspring/decimal"
)

const parameters = `<?xml version="1.0" encoding="iso-8859-1"?>
	<!DOCTYPE XML-FSCM-INVOICE-2003A SYSTEM "XML-FSCM-INVOICE-2003A.DTD"> 
	`

//GetPOSOrders return the las 10 orders from POS
func GetPOSOrders(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	options := struct {
		Limit      string `url:"limit,omitempty,comma"`
		Status     string `url:"status,omitempty,comma"`
		SourceName string `url:"source_name,omitempty,comma"`
		Tag        string `url:"tag,omitempty,comma"`
	}{"10", "any", "pos", "RECHNUNG_DENNER"}

	orders, err := configs.GetClientShop().Order.List(options)
	if err != nil {
		fmt.Println(err)
	}
	options.Tag = "GUTSCHRIFT_DENNER"
	or, err2 := configs.GetClientShop().Order.List(options)
	if err != nil {
		fmt.Println(err2)
	}

	for _, v := range or {
		orders = append(orders, v)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderNumber > orders[j].OrderNumber
	})

	o := calculatePriceEPVariant(orders)
	po := getDateToSentAS2(o)
	models.SendData(w, po)
}

//SendMailPOSOrders from arr of id_sopify send to
func SendMailPOSOrders(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	var results map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &results); err != nil {
		log.Fatal(err)
	}
	orderIDs := fmt.Sprintf("%v", results["order"])
	ordersArr := strings.Split(orderIDs, ",")
	fmt.Println("estas son las ordenes que se van a procesar", ordersArr)

	for _, orderID := range ordersArr {
		xml, err := CreateDennerXML(orderID) //models.SendData(w, orders)

		fmt.Println("el nombre del XML es: " + xml)
		if err != nil {
			fmt.Println(err)
			results["No"] = "las ordenes selecionas no tienen productos para el proveedor  "
		} else {
			//mandar el csv adjunto en un correo
			configs.SendMailForWermProvider(xml)
			configs.CopyFileToAS2(xml)

			err = SaveDateToSenderToAS2(orderID) //save the date to sender a order to AS2 server
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	/* 	xml, err := CreateDennerXML(orderID) //models.SendData(w, orders)

	   	fmt.Println("el nombre del XML es: " + xml)
	   	if err != nil {
	   		fmt.Println(err)
	   		results["No"] = "las ordenes selecionas no tienen productos para el proveedor  "
	   	} else {
	   		//mandar el csv adjunto en un correo
	   		go configs.SendMailForWermProvider(xml)
	   		//go configs.CopyFileToAS2(xml)

	   		err = SaveDateToSenderToAS2(orderID)
	   		if err != nil {
	   			fmt.Println(err)
	   		}
	   	}
	*/
	models.SendData(w, results)
}

//CreateDennerXML create a xml from an id from pos order
func CreateDennerXML(id string) (string, error) {

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return "", err
	}
	o, err2 := configs.GetClientShop().Order.Get(i, nil)
	if err2 != nil {
		return "", err2
	}
	structData := models.NewDennerInvoice()
	dt := time.Now()

	structData.Interchange.IcRef = strconv.Itoa(o.OrderNumber)
	structData.Invoice.Header.MessageReference.ReferenceDate.Date.Date = dt.Format("20060102")
	tags := strings.Split(o.Tags, ",")
	for _, t := range tags {
		if t == "RECHNUNG_DENNER" {
			structData.Invoice.Type = "EFD"
			structData.Invoice.Header.MessageReference.ReferenceDate.ReferenceNo = decimal.NewFromFloat(float64(o.OrderNumber) + 0.1).String()
		} else if t == "GUTSCHRIFT_DENNER" {
			structData.Invoice.Type = "EGS"
			structData.Invoice.Header.MessageReference.ReferenceDate.ReferenceNo = decimal.NewFromFloat(float64(o.OrderNumber) + 0.7).String()
		}
	}
	structData.Invoice.Header.PrintDate.Date.Date = dt.Format("20060102")
	structData.Invoice.Header.DeliverydDate.Date.Date = dt.Format("20060102")
	structData.Invoice.Header.Reference.InvoiceReference.ReferenceDate.ReferenceNo = strconv.Itoa(o.OrderNumber)
	structData.Invoice.Header.Reference.InvoiceReference.ReferenceDate.Date.Date = dt.Format("20060102")
	structData.Invoice.Header.Reference.Order.ReferenceDate.ReferenceNo = strconv.Itoa(o.OrderNumber)
	structData.Invoice.Header.Reference.DeliveryNote.ReferenceDate.ReferenceNo = dt.Format("20060102") + strconv.Itoa(o.OrderNumber) + "360"
	structData.Invoice.Header.Reference.OtherReference[0].ReferenceDate.ReferenceNo = dt.Format("20060102") + strconv.Itoa(o.OrderNumber) + "360"
	structData.Invoice.Header.Biller.DocReferenc.DocReferenc = models.Xvalue + strconv.Itoa(o.OrderNumber) + dt.Format("02")
	structData.Invoice.Header.DeliveryParty.Ean = strings.Split(o.Customer.Note, " ")[1]
	structData.Invoice.Header.DeliveryParty.NaneAddress.Name = o.Customer.DefaultAddress.Company
	structData.Invoice.Header.DeliveryParty.NaneAddress.Street = o.Customer.DefaultAddress.Address1
	structData.Invoice.Header.DeliveryParty.NaneAddress.City = o.Customer.DefaultAddress.City
	structData.Invoice.Header.DeliveryParty.NaneAddress.Zip = o.Customer.DefaultAddress.Zip

	var tAmount, tTax float64

	structData.Invoice.LineItem = make([]models.LineItem, len(o.LineItems), len(o.LineItems))
	for i, p := range o.LineItems {
		structData.Invoice.LineItem[i].LineNumber = strconv.Itoa(i + 1)
		structData.Invoice.LineItem[i].ItemID = make([]models.ItemID, 3, 3)
		structData.Invoice.LineItem[i].ItemID[0].Type = "SA"
		structData.Invoice.LineItem[i].ItemID[0].Data = p.SKU
		structData.Invoice.LineItem[i].ItemID[1].Type = "IN"
		structData.Invoice.LineItem[i].ItemID[1].Data = p.SKU
		structData.Invoice.LineItem[i].ItemID[2].Type = "EN"

		pro, err2 := configs.GetClientShop().Product.Get(p.ProductID, nil)
		fmt.Println(pro)
		if err2 != nil || pro == nil {
			fmt.Println(err2)
			continue
		}
		var variantPrice *decimal.Decimal
		var variantRate *decimal.Decimal

		for _, v := range pro.Variants {
			if v.ID == p.VariantID {
				structData.Invoice.LineItem[i].ItemID[2].Data = v.Barcode
			}
			if v.Title == "EP" {
				variantPrice = v.Price
				if len(p.TaxLines) > 0 {
					variantRate = p.TaxLines[0].Rate
				} else {
					r, _ := decimal.NewFromString("0.077")
					variantRate = &r
				}
			}
		}

		structData.Invoice.LineItem[i].ItemTypeCode = "101"
		structData.Invoice.LineItem[i].ProductName = p.Name

		structData.Invoice.LineItem[i].ItemReference.Type = "ON"
		structData.Invoice.LineItem[i].ItemReference.ReferenceNo = strconv.Itoa(o.OrderNumber)
		structData.Invoice.LineItem[i].ItemReference.LineNo = strconv.Itoa(i + 1)

		structData.Invoice.LineItem[i].Quantity.Type = "47"
		structData.Invoice.LineItem[i].Quantity.Units = "PCE"
		quantyS := calculateQuantity(p.Quantity, p.VariantTitle)
		structData.Invoice.LineItem[i].Quantity.Data = quantyS

		structData.Invoice.LineItem[i].Price.Type = "YYY"
		structData.Invoice.LineItem[i].Price.Units = "PCE"
		structData.Invoice.LineItem[i].Price.Data = variantPrice.String()

		structData.Invoice.LineItem[i].ItemAmount.Type = "66"
		structData.Invoice.LineItem[i].ItemAmount.Amount.Currency = models.Currency
		q, _ := decimal.NewFromString(quantyS)
		ad := variantPrice.Mul(q)

		temp, _ := ad.Float64()
		tAmount += temp
		structData.Invoice.LineItem[i].ItemAmount.Amount.Data = ad.String()

		structData.Invoice.LineItem[i].Tax.TaxBasis.Amount.Currency = models.Currency
		structData.Invoice.LineItem[i].Tax.TaxBasis.Amount.Data = ad.String()

		structData.Invoice.LineItem[i].Tax.Rate.Data = variantRate.Mul(decimal.NewFromInt(100)).String()
		structData.Invoice.LineItem[i].Tax.Rate.Category = models.RateCategory

		tAum := ad.Mul(*variantRate)
		temp2, _ := tAum.Float64()
		tTax += temp2
		structData.Invoice.LineItem[i].Tax.Amount.Data = tAum.Round(2).String()
		structData.Invoice.LineItem[i].Tax.Amount.Currency = models.Currency

	}
	tt := tAmount + tTax
	structData.Invoice.Summary.InvoiceAmount.Amount.Data = fmt.Sprintf("%.2f", tt)
	structData.Invoice.Summary.VatAmount.Amount.Data = fmt.Sprintf("%.2f", tTax)
	//extD := tAmount - tTax
	structData.Invoice.Summary.ExtendedAmount.Amount.Data = fmt.Sprintf("%.2f", tAmount)
	structData.Invoice.Summary.Tax.TaxBasis.Amount.Data = fmt.Sprintf("%.2f", tAmount)
	structData.Invoice.Summary.Tax.Rate.Data = "7.7"
	structData.Invoice.Summary.Tax.Amount.Data = fmt.Sprintf("%.2f", tTax)

	//fmt.Println(structData)
	file, err := xml.MarshalIndent(structData, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	file = []byte(parameters + string(file))
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	sOrderN := strconv.Itoa(o.OrderNumber)
	nameFile := path + "/files/" + sOrderN + "_VICTORY.xml"
	_ = ioutil.WriteFile(nameFile, file, 0644)
	fmt.Println("Created sussefull")
	return nameFile, nil

}

//SendMailPOSRefoundOrders from arr of id_sopify send to
func SendMailPOSRefoundOrders(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	var results map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &results); err != nil {
		log.Fatal(err)
	}
	orderID := fmt.Sprintf("%v", results["order"])
	fmt.Println(string(orderID))

	xml, err := CreateDennerXMLRefound(orderID) //models.SendData(w, orders)

	fmt.Println("el nombre del XML es: " + xml)
	if err != nil {
		//log.Panicln(err)
		results["No"] = "las ordenes selecionas no tienen productos para el proveedor  "
	} else {
		//mandar el csv adjunto en un correo
		go configs.SendMailForWermProvider(xml)
		go configs.CopyFileToAS2(xml)
	}

	models.SendData(w, results)
}

//CreateDennerXMLRefound gd
func CreateDennerXMLRefound(id string) (string, error) {
	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return "", err
	}
	o, err2 := configs.GetClientShop().Order.Get(orderID, nil)
	if err2 != nil {
		return "", err2
	}
	structData := models.NewDennerInvoice()
	dt := time.Now()
	lastR := len(o.Refunds) - 1

	var item []int64
	for _, r := range o.Refunds[lastR].RefundLineItems {
		item = append(item, r.LineItemId)
	}

	structData.Interchange.IcRef = strconv.Itoa(o.OrderNumber)
	structData.Invoice.Type = "EGS"
	structData.Invoice.Header.MessageReference.ReferenceDate.Date.Date = dt.Format("20060102")
	structData.Invoice.Header.MessageReference.ReferenceDate.ReferenceNo = decimal.NewFromFloat(float64(o.OrderNumber) + 0.7).String()
	structData.Invoice.Header.PrintDate.Date.Date = dt.Format("20060102")
	structData.Invoice.Header.DeliverydDate.Date.Date = dt.Format("20060102")
	structData.Invoice.Header.Reference.InvoiceReference.ReferenceDate.ReferenceNo = strconv.Itoa(o.OrderNumber)
	structData.Invoice.Header.Reference.InvoiceReference.ReferenceDate.Date.Date = dt.Format("20060102")
	structData.Invoice.Header.Reference.Order.ReferenceDate.ReferenceNo = strconv.Itoa(o.OrderNumber)
	structData.Invoice.Header.Reference.DeliveryNote.ReferenceDate.ReferenceNo = dt.Format("20060102") + strconv.Itoa(o.OrderNumber) + "360"
	structData.Invoice.Header.Reference.OtherReference[0].ReferenceDate.ReferenceNo = dt.Format("20060102") + strconv.Itoa(o.OrderNumber) + "360"
	structData.Invoice.Header.Biller.DocReferenc.DocReferenc = models.Xvalue + strconv.Itoa(o.OrderNumber) + dt.Format("02")
	structData.Invoice.Header.DeliveryParty.Ean = strings.Split(o.Customer.Note, " ")[1]
	structData.Invoice.Header.DeliveryParty.NaneAddress.Name = o.Customer.DefaultAddress.Company
	structData.Invoice.Header.DeliveryParty.NaneAddress.Street = o.Customer.DefaultAddress.Address1
	structData.Invoice.Header.DeliveryParty.NaneAddress.City = o.Customer.DefaultAddress.City
	structData.Invoice.Header.DeliveryParty.NaneAddress.Zip = o.Customer.DefaultAddress.Zip

	var tAmount, tTax float64

	structData.Invoice.LineItem = make([]models.LineItem, len(item), len(item))
	i := 0
	for _, p := range o.LineItems {
		if _, err := Find(item, o.ID); !err {
			continue
		}
		structData.Invoice.LineItem[i].LineNumber = strconv.Itoa(i + 1)
		structData.Invoice.LineItem[i].ItemID = make([]models.ItemID, 3, 3)
		structData.Invoice.LineItem[i].ItemID[0].Type = "SA"
		structData.Invoice.LineItem[i].ItemID[0].Data = p.SKU
		structData.Invoice.LineItem[i].ItemID[1].Type = "IN"
		structData.Invoice.LineItem[i].ItemID[1].Data = p.SKU
		structData.Invoice.LineItem[i].ItemID[2].Type = "EN"

		pro, err2 := configs.GetClientShop().Product.Get(p.ProductID, nil)
		if err2 != nil {
			fmt.Println(err2)
		}
		var variantPrice *decimal.Decimal
		var variantRate *decimal.Decimal

		for _, v := range pro.Variants {
			if v.ID == p.VariantID {
				structData.Invoice.LineItem[i].ItemID[2].Data = v.Barcode
			}
			if v.Title == "EP" {
				variantPrice = v.Price
				variantRate = p.TaxLines[0].Rate
			}
		}

		structData.Invoice.LineItem[i].ItemTypeCode = "101"
		structData.Invoice.LineItem[i].ProductName = p.Name

		structData.Invoice.LineItem[i].ItemReference.Type = "ON"
		structData.Invoice.LineItem[i].ItemReference.ReferenceNo = strconv.Itoa(o.OrderNumber)
		structData.Invoice.LineItem[i].ItemReference.LineNo = strconv.Itoa(i + 1)

		structData.Invoice.LineItem[i].Quantity.Type = "47"
		structData.Invoice.LineItem[i].Quantity.Units = "PCE"
		quantyS := calculateQuantity(p.Quantity, p.VariantTitle)
		structData.Invoice.LineItem[i].Quantity.Data = quantyS

		structData.Invoice.LineItem[i].Price.Type = "YYY"
		structData.Invoice.LineItem[i].Price.Units = "PCE"
		structData.Invoice.LineItem[i].Price.Data = variantPrice.String()

		structData.Invoice.LineItem[i].ItemAmount.Type = "66"
		structData.Invoice.LineItem[i].ItemAmount.Amount.Currency = models.Currency
		q, _ := decimal.NewFromString(quantyS)
		ad := variantPrice.Mul(q)

		temp, _ := ad.Float64()
		tAmount += temp
		structData.Invoice.LineItem[i].ItemAmount.Amount.Data = ad.String()

		structData.Invoice.LineItem[i].Tax.TaxBasis.Amount.Currency = models.Currency
		structData.Invoice.LineItem[i].Tax.TaxBasis.Amount.Data = ad.String()

		structData.Invoice.LineItem[i].Tax.Rate.Data = variantRate.Mul(decimal.NewFromInt(100)).String()
		structData.Invoice.LineItem[i].Tax.Rate.Category = models.RateCategory

		tAum := ad.Mul(*variantRate)
		temp2, _ := tAum.Float64()
		tTax += temp2
		structData.Invoice.LineItem[i].Tax.Amount.Data = tAum.Round(2).String()
		structData.Invoice.LineItem[i].Tax.Amount.Currency = models.Currency

		i++
	}
	tt := tAmount + tTax
	structData.Invoice.Summary.InvoiceAmount.Amount.Data = fmt.Sprintf("%.2f", tt)
	structData.Invoice.Summary.VatAmount.Amount.Data = fmt.Sprintf("%.2f", tTax)
	//extD := tAmount - tTax
	structData.Invoice.Summary.ExtendedAmount.Amount.Data = fmt.Sprintf("%.2f", tAmount)
	structData.Invoice.Summary.Tax.TaxBasis.Amount.Data = fmt.Sprintf("%.2f", tAmount)
	structData.Invoice.Summary.Tax.Rate.Data = "7.7"
	structData.Invoice.Summary.Tax.Amount.Data = fmt.Sprintf("%.2f", tTax)

	//fmt.Println(structData)
	file, err := xml.MarshalIndent(structData, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	file = []byte(parameters + string(file))
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	sOrderN := strconv.Itoa(o.OrderNumber)
	nameFile := path + "/files/" + sOrderN + "_Refound_VICTORY.xml"
	_ = ioutil.WriteFile(nameFile, file, 0644)
	fmt.Println("Created sussefull")
	return nameFile, nil
}

func SaveDateToSenderToAS2(id string) error {
	query := `INSERT pos_order SET shopify_id = ?, date_send=SYSDATE()`
	result, _ := configs.VicExec(query, id)
	_, err := result.LastInsertId()
	return err
}

func calculateQuantity(units int, amountOfPack string) string {
	unitsPack := strings.Split(amountOfPack, " ")[1]
	//fmt.Println(unitsPack)
	pack, err := strconv.Atoi(unitsPack)
	if err != nil {
		pack = 0
	}
	return fmt.Sprintf("%v", units*pack)
}

//this function calculate a price base on Ep variant
//not used to create XML file
func calculatePriceEPVariant(orders []goshopify.Order) []goshopify.Order {
	resOrder := orders
	//for i, o := range orders {
	for i := 0; i <= 10 && i < len(orders); i++ {

		var tAmount, tTax float64

		for _, p := range orders[i].LineItems {
			pro, err2 := configs.GetClientShop().Product.Get(p.ProductID, nil)
			if err2 != nil {
				fmt.Println(err2)
				continue //some product was delete
			}
			//check if this product has a denner_product tag

			tags := strings.Split(pro.Tags, ",")

			for _, t := range tags {
				if t == "denner-product" {
					resOrder[i].Confirmed = true
					//resOrder[i].FulfillmentStatus = pro.Title
					break
				} else {
					resOrder[i].Confirmed = false
					resOrder[i].FulfillmentStatus = pro.Title
				}
			}

			var variantPrice *decimal.Decimal
			var variantRate *decimal.Decimal

			for _, v := range pro.Variants {
				if v.Title == "EP" {
					variantPrice = v.Price
					if len(p.TaxLines) == 0 {
						ra, _ := decimal.NewFromString("0.077")
						variantRate = &ra
					} else {
						variantRate = p.TaxLines[0].Rate
					}
				}
			}

			quantyS := calculateQuantity(p.Quantity, p.VariantTitle)

			q, _ := decimal.NewFromString(quantyS)
			ad := variantPrice.Mul(q)

			temp, _ := ad.Float64()
			tAmount += temp

			tAum := ad.Mul(*variantRate)
			temp2, _ := tAum.Float64()
			tTax += temp2

		}

		resOrder[i].Reference = fmt.Sprintf("%.2f", tAmount)
	}

	return resOrder
}
func getDateToSentAS2(orders []goshopify.Order) []goshopify.Order {
	resOrder := orders
	for i, o := range orders {
		query := `SELECT date_send FROM pos_order po	
				WHERE po.shopify_id = ?`
		rows, err := configs.VicQuery(query, o.ID)
		if err != nil {
			fmt.Println("error al hacer la consulta para obtener el shopify_id de la DB ", err)
		} else if rows.Next() {
			rows.Scan(&resOrder[i].ProcessedAt)
			resOrder[i].BrowserIp = resOrder[i].ProcessedAt.Format("01-02-2006 15:04")
		} else {
			resOrder[i].BrowserIp = ""
			//err = errors.New("that Oreder no exist")
		}
	}
	return resOrder
}

func Find(slice []int64, val int64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
