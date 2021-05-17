package models

import (
	"log"
	"time"
	"victorydash/configs"
)

type Dok struct {
	DokNr             string    `json:"DokNr,omitempty"`
	DokTyp            string    `json:"DokTyp,omitempty"`
	DokDat            time.Time `json:"DokDat,omitempty"`
	LFirma            string    `json:"LFirma,omitempty"`
	LStrasse          string    `json:"LStrasse,omitempty"`
	LPLZ              string    `json:"LPLZ,omitempty"`
	LOrt              string    `json:"LOrt,omitempty"`
	TotalPos          float64   `json:"TotalPos,omitempty"`
	TotalSteuer       float64   `json:"TotalSteuer,omitempty"`
	Zahlung           string    `json:"Zahlung,omitempty"`
	TotalDok          float64   `json:"TotalDok,omitempty"`
	EBPPBillAccountID string    `json:"EBPPBillAccountID,omitempty"`
	Steuersatz1       float64   `json:"Steuersatz1,omitempty"`

	Products []Product `json:"Products,omitempty"`
}
type Product struct {
	HeadPosNr     int     `json:"HeadPosNr,omitempty"`
	ArtNr         string  `json:"ArtNr,omitempty"`
	BezeichnungD1 string  `json:"BezeichnungD1,omitempty"`
	Bestellmenge  int     `json:"Bestellmenge,omitempty"`
	VP            float64 `json:"VP,omitempty"`
	PosiTot       float64 `json:"PosiTot,omitempty"`
	MWStBtrg      float64 `json:"MWStBtrg,omitempty"`
	IEANCode      string  `json:"iEANCode,omitempty"`
}
type DokScan struct {
	DokNr       string    `json:"DokNr,omitempty"`
	DokTyp      string    `json:"DokTyp,omitempty"`
	DokDat      time.Time `json:"DokDat,omitempty"`
	LFirma      string    `json:"LFirma,omitempty"`
	LStrasse    string    `json:"LStrasse,omitempty"`
	LPLZ        string    `json:"LPLZ,omitempty"`
	LOrt        string    `json:"LOrt,omitempty"`
	TotalPos    float64   `json:"TotalPos,omitempty"`
	TotalSteuer float64   `json:"TotalSteuer,omitempty"`
	Zahlung     string    `json:"Zahlung,omitempty"`
	TotalDok    float64   `json:"TotalDok,omitempty"`
	Steuersatz1 float64   `json:"Steuersatz1,omitempty"`

	HeadPosNr         int     `json:"HeadPosNr,omitempty"`
	ArtNr             string  `json:"ArtNr,omitempty"`
	BezeichnungD1     string  `json:"BezeichnungD1,omitempty"`
	Bestellmenge      int     `json:"Bestellmenge,omitempty"`
	VP                float64 `json:"VP,omitempty"`
	PosiTot           float64 `json:"PosiTot,omitempty"`
	MWStBtrg          float64 `json:"MWStBtrg,omitempty"`
	IEANCode          string  `json:"iEANCode,omitempty"`
	EBPPBillAccountID string  `json:"EBPPBillAccountID,omitempty"`
}

func GetDokByDokNur(dokNur string) (Dok, error) {
	query := `SELECT DokNr, DokTyp, DokDat, LFirma, LStrasse, LPLZ, LOrt, TotalPos,        
	TotalSteuer, Zahlung, TotalDok, HeadPosNr, ArtNr, BezeichnungD1,    
	Bestellmenge, VP, PosiTot, MWStBtrg, IEANCode, EBPPBillAccountID, Steuersatz1
	from AUFTRAG_VS_TEST.dbo.X_Dok_Denner WHERE DokNr = @p1 `

	rows, err := configs.SageQuery(query, dokNur)

	if err != nil {
		log.Fatal(err)
	}
	var dokN string
	dok := Dok{}

	for rows.Next() {

		dokScan := DokScan{}
		err := rows.Scan(&dokScan.DokNr, &dokScan.DokTyp, &dokScan.DokDat, &dokScan.LFirma, &dokScan.LStrasse,
			&dokScan.LPLZ, &dokScan.LOrt, &dokScan.TotalPos, &dokScan.TotalSteuer, &dokScan.Zahlung,
			&dokScan.TotalDok, &dokScan.HeadPosNr, &dokScan.ArtNr, &dokScan.BezeichnungD1, &dokScan.Bestellmenge,
			&dokScan.VP, &dokScan.PosiTot, &dokScan.MWStBtrg, &dokScan.IEANCode, &dokScan.EBPPBillAccountID,
			&dokScan.Steuersatz1)
		if err != nil {
			log.Fatal(err)
		}

		if dokN == dokScan.DokNr {
			p := Product{}
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

			dokN = dokScan.DokNr
			dok = Dok{}

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
			dok.Steuersatz1 = dokScan.Steuersatz1

			dok.Products = make([]Product, 0)
			p := Product{}
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

	return dok, err
}
