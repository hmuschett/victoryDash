package models

import (
	"encoding/xml"
)

type DennerInvoice struct {
	XMLName     xml.Name    `xml:"XML-FSCM-INVOICE-2003A"`
	Interchange Interchange `xml:"INTERCHANGE"`
	Invoice     Invoice     `xml:"INVOICE"`
}
type Interchange struct {
	XMLName   xml.Name  `xml:"INTERCHANGE"`
	IcSender  IcSender  `xml:"IC-SENDER"`
	IcReciver IcReciver `xml:"IC-RECEIVER"`
	IcRef     int64     `xml:"IC-Ref"`
}
type IcSender struct {
	XMLName xml.Name `xml:"IC-SENDER"`
	Pid     string   `xml:"Pid"`
}
type IcReciver struct {
	XMLName xml.Name `xml:"IC-RECEIVER"`
	Pid     string   `xml:"Pid"`
}
type Invoice struct {
	XMLName  xml.Name   `xml:"INVOICE"`
	Type     string     `xml:"Type,attr"`
	Header   Header     `xml:"HEADER"`
	LineItem []LineItem `xml:"LINE-ITEM"`
	Summary  Summary    `xml:"SUMMARY"`
}
type Header struct {
	XMLName          xml.Name         `xml:"HEADER"`
	MessageReference MessageReference `xml:"MESSAGE-REFERENCE"`
	PrintDate        PrintDate        `xml:"PRINT-DATE"`
	DeliverydDate    DeliverydDate    `xml:"DELIVERY-DATE"`
	Reference        Reference        `xml:"REFERENCE"`
	Biller           Biller           `xml:"BILLER"`
	Payer            Payer            `xml:"PAYER"`
	DeliveryParty    DeliveryParty    `xml:"DELIVERY-PARTY"`
}
type MessageReference struct {
	XMLName       xml.Name      `xml:"MESSAGE-REFERENCE"`
	ReferenceDate ReferenceDate `xml:"REFERENCE-DATE"`
}
type Date struct {
	XMLName xml.Name `xml:"Date,omitempty"`
	Date    string   `xml:",chardata"`
	Format  string   `xml:"Format,attr,omitempty"`
}
type ReferenceDate struct {
	XMLName     xml.Name `xml:"REFERENCE-DATE"`
	ReferenceNo string   `xml:"Reference-No"`
	Date        *Date
}
type PrintDate struct {
	XMLName xml.Name `xml:"PRINT-DATE"`
	Date    Date     `xml:"Date"`
}
type DeliverydDate struct {
	XMLName xml.Name `xml:"DELIVERY-DATE"`
	Date    Date     `xml:"Date"`
}
type Reference struct {
	XMLName          xml.Name         `xml:"REFERENCE"`
	InvoiceReference InvoiceReference `xml:"INVOICE-REFERENCE"`
	Order            Order            `xml:"ORDER"`
	DeliveryNote     DeliveryNote     `xml:"DELIVERY-NOTE"`
	OtherReference   []OtherReference `xml:"OTHER-REFERENCE"`
}
type InvoiceReference struct {
	XMLName       xml.Name      `xml:"INVOICE-REFERENCE"`
	ReferenceDate ReferenceDate `xml:"REFERENCE-DATE"`
}
type Order struct {
	XMLName       xml.Name      `xml:"ORDER"`
	ReferenceDate ReferenceDate `xml:"REFERENCE-DATE"`
}
type DeliveryNote struct {
	XMLName       xml.Name      `xml:"DELIVERY-NOTE"`
	ReferenceDate ReferenceDate `xml:"REFERENCE-DATE"`
}
type OtherReference struct {
	XMLName       xml.Name      `xml:"OTHER-REFERENCE"`
	ReferenceDate ReferenceDate `xml:"REFERENCE-DATE"`
	Type          string        `xml:"Type,attr"`
}
type Biller struct {
	XMLName     xml.Name `xml:"BILLER"`
	TaxNo       string   `xml:"Tax-No"`
	DocReferenc DocReferenc
	Ean         string `xml:"PARTY-ID>Ean"`
	NaneAddress NaneAddress
	BankInfo    BankInfo
}
type DocReferenc struct {
	XMLName     xml.Name `xml:"Doc-Reference"`
	DocReferenc string   `xml:",chardata"`
	Type        string   `xml:"Type,attr"`
}
type NaneAddress struct {
	XMLName xml.Name `xml:"NAME-ADDRESS"`
	Name    string   `xml:"NAME>Line-35"`
	Street  string   `xml:"STREET>Line-35"`
	City    string   `xml:"City"`
	Zip     string   `xml:"Zip"`
	Country string   `xml:"Country"`
}
type BankInfo struct {
	XMLName xml.Name `xml:"BANK-INFO"`
	AcctNo  string   `xml:"Acct-No"`
	BankId  BankId
}
type BankId struct {
	XMLName xml.Name `xml:"BankId"`
	BankID  string   `xml:",chardata"`
	Country string   `xml:"Country,attr"`
	Type    string   `xml:"Type,attr"`
}
type Payer struct {
	XMLName     xml.Name `xml:"PAYER"`
	TaxNo       string   `xml:"Tax-No"`
	Ean         string   `xml:"PARTY-ID>Ean"`
	NaneAddress NaneAddress
}
type DeliveryParty struct {
	XMLName     xml.Name `xml:"DELIVERY-PARTY"`
	Ean         string   `xml:"Ean"`
	NaneAddress NaneAddress
}
type LineItem struct {
	XMLName       xml.Name      `xml:"LINE-ITEM"`
	LineNumber    string        `xml:"Line-Number,attr"`
	ItemID        []ItemID      `xml:"ITEM-ID>Item-Id"`
	ItemTypeCode  string        `xml:"ITEM-DESCRIPTION>Item-Type-Code"`
	ProductName   string        `xml:"ITEM-DESCRIPTION>Line-35"`
	ItemReference ItemReference `xml:"ITEM-REFERENCE"`
	Quantity      Quantity      `xml:"Quantity"`
	Price         Price         `xml:"Price"`
	ItemAmount    ItemAmount    `xml:"ITEM-AMOUNT"`
	Tax           Tax           `xml:"TAX"`
}
type ItemID struct {
	XMLName xml.Name `xml:"Item-Id"`
	Data    string   `xml:",chardata"`
	Type    string   `xml:"Type,attr"`
}
type ItemReference struct {
	XMLName     xml.Name `xml:"ITEM-REFERENCE"`
	Type        string   `xml:"Type,attr"`
	ReferenceNo string   `xml:"REFERENCE-DATE>Reference-No"`
	LineNo      string   `xml:"REFERENCE-DATE>Line-No"`
}
type Quantity struct {
	XMLName xml.Name `xml:"Quantity"`
	Data    string   `xml:",chardata"`
	Type    string   `xml:"Type,attr"`
	Units   string   `xml:"Units,attr"`
}
type Price struct {
	XMLName xml.Name `xml:"Price"`
	Data    string   `xml:",chardata"`
	Type    string   `xml:"Type,attr"`
	Units   string   `xml:"Units,attr"`
}
type ItemAmount struct {
	XMLName xml.Name `xml:"ITEM-AMOUNT"`
	Type    string   `xml:"Type,attr"`
	Amount  Amount   `xml:"Amount"`
}
type Amount struct {
	//XMLName  xml.Name `xml:"AMOUNT"`
	Data     string `xml:",chardata"`
	Currency string `xml:"Currency,attr"`
}
type Tax struct {
	XMLName  xml.Name `xml:"TAX"`
	TaxBasis TaxBasis `xml:"TAX-BASIS"`
	Rate     Rate     `xml:"Rate"`
	Amount   Amount   `xml:"Amount"`
}
type TaxBasis struct {
	XMLName xml.Name `xml:"TAX-BASIS"`
	Amount  Amount   `xml:"Amount"`
}
type Rate struct {
	XMLName  xml.Name `xml:"Rate"`
	Data     string   `xml:",chardata"`
	Category string   `xml:"Category,attr"`
}
type Summary struct {
	XMLName        xml.Name       `xml:"SUMMARY"`
	InvoiceAmount  InvoiceAmount  `xml:"INVOICE-AMOUNT"`
	VatAmount      VatAmount      `xml:"VAT-AMOUNT"`
	ExtendedAmount ExtendedAmount `xml:"EXTENDED-AMOUNT"`
	Tax            Tax            `xml:"TAX"`
	PaymentTerms   PaymentTerms   `xml:"PAYMENT-TERMS"`
}
type InvoiceAmount struct {
	XMLName     xml.Name `xml:"INVOICE-AMOUNT"`
	PrintStatus string   `xml:"Print-Status,attr"`
	Amount      Amount   `xml:"Amount"`
}
type VatAmount struct {
	XMLName     xml.Name `xml:"VAT-AMOUNT"`
	PrintStatus string   `xml:"Print-Status,attr"`
	Amount      Amount   `xml:"Amount"`
}
type ExtendedAmount struct {
	XMLName xml.Name `xml:"EXTENDED-AMOUNT"`
	Type    string   `xml:"Type,attr"`
	Amount  Amount   `xml:"Amount"`
}
type PaymentTerms struct {
	XMLName xml.Name `xml:"PAYMENT-TERMS"`
	Basic   Basic    `xml:"BASIC"`
}
type Basic struct {
	XMLName       xml.Name      `xml:"BASIC"`
	TermsType     string        `xml:"Terms-Type,attr"`
	PaymentType   string        `xml:"Payment-Type,attr"`
	PaymentPeriod PaymentPeriod `xml:"TERMS>Payment-Period"`
}
type PaymentPeriod struct {
	XMLName      xml.Name `xml:"Payment-Period"`
	ReferenceDay string   `xml:"Reference-Day,attr"`
	Type         string   `xml:"type,attr"`
	OnOrAfter    string   `xml:"On-Or-After,attr"`
	Data         string   `xml:",chardata"`
}

const (
	pidSender          = "7640146250001"
	pidReciver         = "7610029000009"
	icRef              = 89142
	typeEFD            = "EFD"
	formatDate         = "CCYYMMDD"
	typeCR             = "CR"
	typeIT             = "IT"
	typeESRNEU         = "ESR-NEU"
	taxNo              = "CHE274334848"
	eamBiller          = "7640146250001"
	nameVictory        = "Victory Switzerland GmbH"
	streetBiller       = "Solothurnstrasse 24 C"
	cityBiller         = "Kirchberg"
	zipBiller          = "3422"
	country            = "CH"
	acctNo             = "010001456"
	typeBCNr           = "BCNr-nat"
	backID             = "001996"
	eamPayer           = "7610029000009" //es el moismo valor que pidReciver
	namePayer          = "Denner AG"
	streetPayer        = "Grubenstrasse 10"
	cityPayer          = "Zürich"
	zipPayer           = "8045"
	itemTypeCode       = "101"
	itNo               = "400000"
	Xvalue             = "26638200400000000000"
	Currency           = "CHF"
	printStatus        = "25"
	typeExtendedAmount = "25"
	RateCategory       = "S"
	typeBasicTerms     = "1"
	typePayment        = "ESR"
	referenceDay       = "5"
	typePaymentPeriod  = "CD"
	onOrAfter          = "3"
	paymentPeriod      = "10"
)

func NewDennerInvoice() DennerInvoice {
	data := DennerInvoice{}
	data.Interchange.IcSender.Pid = pidSender
	data.Interchange.IcReciver.Pid = pidReciver
	data.Interchange.IcRef = icRef
	data.Invoice.Type = typeEFD
	data.Invoice.Header.MessageReference.ReferenceDate.Date = new(Date)
	data.Invoice.Header.MessageReference.ReferenceDate.Date.Format = formatDate
	data.Invoice.Header.MessageReference.ReferenceDate.Date.Date = "37373737"
	data.Invoice.Header.PrintDate.Date.Format = formatDate
	data.Invoice.Header.DeliverydDate.Date.Format = formatDate
	data.Invoice.Header.Reference.InvoiceReference.ReferenceDate.Date = new(Date)
	data.Invoice.Header.Reference.InvoiceReference.ReferenceDate.Date.Format = formatDate
	data.Invoice.Header.Reference.InvoiceReference.ReferenceDate.Date.Date = "65655656"
	data.Invoice.Header.Reference.OtherReference = make([]OtherReference, 2, 2)
	data.Invoice.Header.Reference.OtherReference[0].Type = typeCR
	data.Invoice.Header.Reference.OtherReference[1].Type = typeIT
	data.Invoice.Header.Reference.OtherReference[1].ReferenceDate.ReferenceNo = itNo

	data.Invoice.Header.Biller.TaxNo = taxNo
	data.Invoice.Header.Biller.DocReferenc.Type = typeESRNEU
	data.Invoice.Header.Biller.Ean = eamBiller
	data.Invoice.Header.Biller.NaneAddress.Name = nameVictory
	data.Invoice.Header.Biller.NaneAddress.Street = streetBiller
	data.Invoice.Header.Biller.NaneAddress.City = cityBiller
	data.Invoice.Header.Biller.NaneAddress.Zip = zipBiller
	data.Invoice.Header.Biller.NaneAddress.Country = country

	data.Invoice.Header.Biller.BankInfo.AcctNo = acctNo
	data.Invoice.Header.Biller.BankInfo.BankId.Country = country
	data.Invoice.Header.Biller.BankInfo.BankId.Type = typeBCNr
	data.Invoice.Header.Biller.BankInfo.BankId.BankID = backID

	data.Invoice.Header.Payer.TaxNo = taxNo
	data.Invoice.Header.Payer.Ean = eamPayer
	data.Invoice.Header.Payer.NaneAddress.Name = namePayer
	data.Invoice.Header.Payer.NaneAddress.Street = streetPayer
	data.Invoice.Header.Payer.NaneAddress.City = cityPayer
	data.Invoice.Header.Payer.NaneAddress.Zip = zipPayer
	data.Invoice.Header.Payer.NaneAddress.Country = country

	data.Invoice.Summary.InvoiceAmount.PrintStatus = printStatus
	data.Invoice.Summary.InvoiceAmount.Amount.Currency = Currency
	data.Invoice.Summary.VatAmount.PrintStatus = printStatus
	data.Invoice.Summary.VatAmount.Amount.Currency = Currency
	data.Invoice.Summary.ExtendedAmount.Type = typeExtendedAmount
	data.Invoice.Summary.ExtendedAmount.Amount.Currency = Currency
	data.Invoice.Summary.Tax.TaxBasis.Amount.Currency = Currency
	data.Invoice.Summary.Tax.Rate.Category = RateCategory
	data.Invoice.Summary.Tax.Amount.Currency = Currency
	data.Invoice.Summary.PaymentTerms.Basic.TermsType = typeBasicTerms
	data.Invoice.Summary.PaymentTerms.Basic.PaymentType = typePayment

	data.Invoice.Summary.PaymentTerms.Basic.PaymentPeriod.ReferenceDay = referenceDay
	data.Invoice.Summary.PaymentTerms.Basic.PaymentPeriod.Type = typePaymentPeriod
	data.Invoice.Summary.PaymentTerms.Basic.PaymentPeriod.OnOrAfter = onOrAfter
	data.Invoice.Summary.PaymentTerms.Basic.PaymentPeriod.Data = paymentPeriod

	return data
}
