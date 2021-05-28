package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"victorydash/configs"
	"victorydash/models"
	"victorydash/utils"
)

func GetSageOrdersByDates(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	fromDate := r.FormValue("fromDate")
	toDate := r.FormValue("toDate")

	doks := make([]models.Dok, 0)

	query := `SELECT DokNr, DokTyp, DokDat, LFirma, LStrasse, LPLZ, LOrt, TotalPos,        
				TotalSteuer, Zahlung, TotalDok, HeadPosNr, ArtNr, BezeichnungD1,    
				Bestellmenge, VP, PosiTot, MWStBtrg, IEANCode, EBPPBillAccountID
				from AUFTRAG_VS_TEST.dbo.X_Dok_Denner
				WHERE DokDat BETWEEN  @p1  and @p2 ORDER BY DokDat DESC`

	fromD, _ := time.Parse("2006-1-2", fromDate)
	toD, _ := time.Parse("2006-1-2", toDate)
	rows, err := configs.SageQuery(query, fromD, toD)

	if err != nil {
		log.Fatal(err)
	}
	var dokN string
	dok := models.Dok{}

	for rows.Next() {

		dokScan := models.DokScan{}
		err := rows.Scan(&dokScan.DokNr, &dokScan.DokTyp, &dokScan.DokDat, &dokScan.LFirma, &dokScan.LStrasse,
			&dokScan.LPLZ, &dokScan.LOrt, &dokScan.TotalPos, &dokScan.TotalSteuer, &dokScan.Zahlung,
			&dokScan.TotalDok, &dokScan.HeadPosNr, &dokScan.ArtNr, &dokScan.BezeichnungD1, &dokScan.Bestellmenge,
			&dokScan.VP, &dokScan.PosiTot, &dokScan.MWStBtrg, &dokScan.IEANCode, &dokScan.EBPPBillAccountID)
		if err != nil {
			log.Fatal(err)
		}

		if dokN == dokScan.DokNr {
			p := models.Product{}
			p.HeadPosNr = dokScan.HeadPosNr
			p.ArtNr = dokScan.ArtNr
			p.BezeichnungD1 = dokScan.BezeichnungD1
			p.Bestellmenge = dokScan.Bestellmenge
			p.VP = dokScan.VP
			p.PosiTot = dokScan.PosiTot
			p.MWStBtrg = dokScan.MWStBtrg
			p.IEANCode = dokScan.IEANCode

			dok.Products = append(dok.Products, p)
		} else {
			if dokN != "" {
				doks = append(doks, dok)
			}

			dokN = dokScan.DokNr
			dok = models.Dok{}

			dok.DokNr = dokScan.DokNr
			dok.DokTyp = dokScan.DokTyp
			dok.DokDat = dokScan.DokDat
			dok.LFirma = dokScan.LFirma
			dok.LStrasse = dokScan.LStrasse
			dok.LPLZ = dokScan.LPLZ
			dok.LOrt = dokScan.LOrt
			dok.TotalPos = dokScan.TotalPos
			dok.TotalSteuer = dokScan.TotalSteuer
			dok.Zahlung = "10"
			dok.TotalDok = dokScan.TotalDok
			dok.EBPPBillAccountID = dokScan.EBPPBillAccountID

			dok.Products = make([]models.Product, 0)
			p := models.Product{}
			p.HeadPosNr = dokScan.HeadPosNr
			p.ArtNr = dokScan.ArtNr
			p.BezeichnungD1 = dokScan.BezeichnungD1
			p.Bestellmenge = dokScan.Bestellmenge
			p.VP = dokScan.VP
			p.PosiTot = dokScan.PosiTot
			p.MWStBtrg = dokScan.MWStBtrg
			p.IEANCode = dokScan.IEANCode

			dok.Products = append(dok.Products, p)
		}
	}
	doks = append(doks, dok) //añadir al ultimo dok
	d := getDateToSentAS2Dok(doks)
	models.SendData(w, d)
}

//GetSageOrders return the las 10 orders from POS
func GetSageOrders(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w)
	doks := make([]models.Dok, 0)

	query := `SELECT DokNr, DokTyp, DokDat, LFirma, LStrasse, LPLZ, LOrt, TotalPos,        
						TotalSteuer, Zahlung, TotalDok, HeadPosNr, ArtNr, BezeichnungD1,    
						Bestellmenge, VP, PosiTot, MWStBtrg, IEANCode, EBPPBillAccountID
						from AUFTRAG_VS_TEST.dbo.X_Dok_Denner
							WHERE DokNr in (SELECT  dok.nn from (
							SELECT  DISTINCT  x.DokNr as nn , x.DokDat  FROM dbo.X_Dok_Denner x    
							    ORDER BY x.DokDat DESC  
							    OFFSET 0 ROWS
							    FETCH NEXT 20 ROWS ONLY) as dok) `
	rows, err := configs.SageQuery(query)

	if err != nil {
		log.Fatal(err)
	}
	var dokN string
	dok := models.Dok{}

	for rows.Next() {

		dokScan := models.DokScan{}
		err := rows.Scan(&dokScan.DokNr, &dokScan.DokTyp, &dokScan.DokDat, &dokScan.LFirma, &dokScan.LStrasse,
			&dokScan.LPLZ, &dokScan.LOrt, &dokScan.TotalPos, &dokScan.TotalSteuer, &dokScan.Zahlung,
			&dokScan.TotalDok, &dokScan.HeadPosNr, &dokScan.ArtNr, &dokScan.BezeichnungD1, &dokScan.Bestellmenge,
			&dokScan.VP, &dokScan.PosiTot, &dokScan.MWStBtrg, &dokScan.IEANCode, &dokScan.EBPPBillAccountID)
		if err != nil {
			log.Fatal(err)
		}

		if dokN == dokScan.DokNr {
			p := models.Product{}
			p.HeadPosNr = dokScan.HeadPosNr
			p.ArtNr = dokScan.ArtNr
			p.BezeichnungD1 = dokScan.BezeichnungD1
			p.Bestellmenge = dokScan.Bestellmenge
			p.VP = dokScan.VP
			p.PosiTot = dokScan.PosiTot
			p.MWStBtrg = dokScan.MWStBtrg
			p.IEANCode = dokScan.IEANCode

			dok.Products = append(dok.Products, p)
		} else {
			if dokN != "" {
				doks = append(doks, dok)
			}

			dokN = dokScan.DokNr
			dok = models.Dok{}

			dok.DokNr = dokScan.DokNr
			dok.DokTyp = dokScan.DokTyp
			dok.DokDat = dokScan.DokDat
			dok.LFirma = dokScan.LFirma
			dok.LStrasse = dokScan.LStrasse
			dok.LPLZ = dokScan.LPLZ
			dok.LOrt = dokScan.LOrt
			dok.TotalPos = dokScan.TotalPos
			dok.TotalSteuer = dokScan.TotalSteuer
			dok.Zahlung = "10"
			dok.TotalDok = dokScan.TotalDok
			dok.EBPPBillAccountID = dokScan.EBPPBillAccountID

			dok.Products = make([]models.Product, 0)
			p := models.Product{}
			p.HeadPosNr = dokScan.HeadPosNr
			p.ArtNr = dokScan.ArtNr
			p.BezeichnungD1 = dokScan.BezeichnungD1
			p.Bestellmenge = dokScan.Bestellmenge
			p.VP = dokScan.VP
			p.PosiTot = dokScan.PosiTot
			p.MWStBtrg = dokScan.MWStBtrg
			p.IEANCode = dokScan.IEANCode

			dok.Products = append(dok.Products, p)
		}
	}
	doks = append(doks, dok) //añadir al ultimo dok
	d := getDateToSentAS2Dok(doks)
	models.SendData(w, d)
}

//SendOrders from arr of id_sopify send to
func SendOrders(w http.ResponseWriter, r *http.Request) {
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
		xml, err := CreateDennerXMLFromSage(orderID) //models.SendData(w, orders)

		fmt.Println("el nombre del XML es: " + xml)
		if err != nil {
			fmt.Println(err)
			results["No"] = "las ordenes selecionas no tienen productos para el proveedor  "
		} else {
			//mandar el csv adjunto en un correo
			configs.SendMailForWermProvider(xml)
			configs.CopyFileToAS2(xml)

			err = saveDateToSenderToAS2Server(orderID) //save the date to sender a order to AS2 server
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	models.SendData(w, results)
}

//CreateDennerXMLFromSage create a xml from an id from pos order
func CreateDennerXMLFromSage(id string) (string, error) {

	o, err2 := models.GetDokByDokNur(id)
	if err2 != nil {
		return "", err2
	}
	structData := models.NewDennerInvoice()

	structData.Interchange.IcRef = o.DokNr
	structData.Invoice.Header.MessageReference.ReferenceDate.Date.Date = o.DokDat.Format("20060102")

	if o.DokTyp == "Rechnung_Denner" {
		structData.Invoice.Type = "EFD"
	} else if o.DokTyp == "Gutschrift_Denner" {
		structData.Invoice.Type = "EGS"
	}
	structData.Invoice.Header.MessageReference.ReferenceDate.ReferenceNo = o.DokNr
	structData.Invoice.Header.PrintDate.Date.Date = o.DokDat.Format("20060102")
	structData.Invoice.Header.DeliverydDate.Date.Date = o.DokDat.Format("20060102")
	structData.Invoice.Header.Reference.InvoiceReference.ReferenceDate.ReferenceNo = o.DokNr
	structData.Invoice.Header.Reference.InvoiceReference.ReferenceDate.Date.Date = o.DokDat.Format("20060102")
	structData.Invoice.Header.Reference.Order.ReferenceDate.ReferenceNo = o.DokNr
	structData.Invoice.Header.Reference.DeliveryNote.ReferenceDate.ReferenceNo = o.DokDat.Format("20060102") + o.DokNr + "360"
	structData.Invoice.Header.Reference.OtherReference[0].ReferenceDate.ReferenceNo = o.DokDat.Format("20060102") + o.DokNr + "360"
	structData.Invoice.Header.Biller.DocReferenc.DocReferenc = models.Xvalue + o.DokNr + o.DokDat.Format("02")
	structData.Invoice.Header.DeliveryParty.Ean = o.EBPPBillAccountID
	structData.Invoice.Header.DeliveryParty.NaneAddress.Name = o.LFirma
	structData.Invoice.Header.DeliveryParty.NaneAddress.Street = o.LStrasse
	structData.Invoice.Header.DeliveryParty.NaneAddress.City = o.LOrt
	structData.Invoice.Header.DeliveryParty.NaneAddress.Zip = o.LPLZ

	structData.Invoice.LineItem = make([]models.LineItem, len(o.Products), len(o.Products))
	for i, p := range o.Products {
		structData.Invoice.LineItem[i].LineNumber = strconv.Itoa(i + 1)
		structData.Invoice.LineItem[i].ItemID = make([]models.ItemID, 3, 3)
		structData.Invoice.LineItem[i].ItemID[0].Type = "SA"
		structData.Invoice.LineItem[i].ItemID[0].Data = p.ArtNr
		structData.Invoice.LineItem[i].ItemID[1].Type = "IN"
		structData.Invoice.LineItem[i].ItemID[1].Data = p.ArtNr
		structData.Invoice.LineItem[i].ItemID[2].Type = "EN"
		structData.Invoice.LineItem[i].ItemID[2].Data = p.IEANCode

		structData.Invoice.LineItem[i].ItemTypeCode = "101"
		structData.Invoice.LineItem[i].ProductName = p.BezeichnungD1

		structData.Invoice.LineItem[i].ItemReference.Type = "ON"
		structData.Invoice.LineItem[i].ItemReference.ReferenceNo = p.ArtNr
		structData.Invoice.LineItem[i].ItemReference.LineNo = strconv.Itoa(i + 1)

		structData.Invoice.LineItem[i].Quantity.Type = "47"
		structData.Invoice.LineItem[i].Quantity.Units = "PCE"
		structData.Invoice.LineItem[i].Quantity.Data = strconv.Itoa(p.Bestellmenge)

		structData.Invoice.LineItem[i].Price.Type = "YYY"
		structData.Invoice.LineItem[i].Price.Units = "PCE"
		structData.Invoice.LineItem[i].Price.Data = fmt.Sprintf("%.2f", p.VP)

		structData.Invoice.LineItem[i].ItemAmount.Type = "66"
		structData.Invoice.LineItem[i].ItemAmount.Amount.Currency = models.Currency
		structData.Invoice.LineItem[i].ItemAmount.Amount.Data = fmt.Sprintf("%.2f", p.PosiTot)

		structData.Invoice.LineItem[i].Tax.TaxBasis.Amount.Currency = models.Currency
		structData.Invoice.LineItem[i].Tax.TaxBasis.Amount.Data = fmt.Sprintf("%.2f", p.PosiTot)

		structData.Invoice.LineItem[i].Tax.Rate.Category = models.RateCategory
		structData.Invoice.LineItem[i].Tax.Rate.Data = fmt.Sprintf("%.2f", o.Steuersatz1)

		structData.Invoice.LineItem[i].Tax.Amount.Currency = models.Currency
		structData.Invoice.LineItem[i].Tax.Amount.Data = fmt.Sprintf("%.2f", p.MWStBtrg)

	}

	structData.Invoice.Summary.InvoiceAmount.Amount.Data = fmt.Sprintf("%.2f", o.TotalDok)
	structData.Invoice.Summary.VatAmount.Amount.Data = fmt.Sprintf("%.2f", o.TotalSteuer)
	structData.Invoice.Summary.ExtendedAmount.Amount.Data = fmt.Sprintf("%.2f", o.TotalPos)
	structData.Invoice.Summary.Tax.TaxBasis.Amount.Data = fmt.Sprintf("%.2f", o.TotalPos)
	structData.Invoice.Summary.Tax.Rate.Data = fmt.Sprintf("%.2f", o.Steuersatz1)
	structData.Invoice.Summary.Tax.Amount.Data = fmt.Sprintf("%.2f", o.TotalSteuer)

	file, err := xml.MarshalIndent(structData, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	file = []byte(parameters + string(file))
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	sOrderN := o.DokNr
	nameFile := path + "/files/" + sOrderN + "_VICTORY.xml"
	_ = ioutil.WriteFile(nameFile, file, 0644)
	fmt.Println("Created sussefull")
	return nameFile, nil
}
func getDateToSentAS2Dok(orders []models.Dok) []models.Dok {
	resOrder := orders

	for i, o := range orders {
		if strings.Trim(o.DokNr, " ") != "" {
			query := `SELECT date_send FROM dok_send_as po	
				WHERE po.shopify_id = ?`
			rows, err := configs.VicQuery(query, o.DokNr)

			if err != nil {
				fmt.Println("error al hacer la consulta para obtener el shopify_id de la DB ", err)
			} else if rows.Next() {
				rows.Scan(&resOrder[i].SentAS)
				resOrder[i].Zahlung = resOrder[i].SentAS.Format("01-02-2006 15:04")

			} else {
				resOrder[i].Zahlung = "hhh"
			}
		}
	}
	return resOrder
}
func saveDateToSenderToAS2Server(id string) error {
	query := `INSERT dok_send_as SET shopify_id = ?, date_send=SYSDATE()`
	result, _ := configs.VicExec(query, id)
	_, err := result.LastInsertId()
	return err
}
