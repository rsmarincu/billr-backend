package domain

type QuantityType string

const (
	QuantityTypeHours QuantityType = "hours"
	QuantityTypeDays  QuantityType = "days"
	QuantityTypeUnit  QuantityType = "unit"
)

type Invoice struct {
	Id              string
	InvoiceNumber   string
	UserCompany     Company
	InvoicedCompany Company
	Currency        string
	Created         string
	Due             string
	Articles        []Article
	Total           float64
}

type Company struct {
	Id                 string
	Name               string
	RegistrationNumber string
	CUI                string
	VatId              string
	Email              string
	BankAccount        BankAccount
	StreetAddress      string
	City               string
	Country            string
	PostCode           string
}

type Article struct {
	Description  string
	Quantity     int
	QuantityType QuantityType
	Price        float64
	Amount       float64
}

type BankAccount struct {
	Name string
	IBAN string
}
