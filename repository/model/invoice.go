package model

import (
	"time"

	"github.com/rsmarincu/billr/domain"
)

type ArticleRecords []ArticleRecord

const layoutStr = "2006-10-02"

type InvoiceRecord struct {
	Id        string    `db:"id"`
	UserId    string    `db:"userId"`
	CompanyId string    `db:"userCompanyId"`
	ClientId  string    `db:"clientId"`
	Currency  string    `db:"currency"`
	Created   time.Time `db:"created"`
	Due       time.Time `db:"due"`
	Total     float64   `db:"total"`
}

type ArticleRecord struct {
	Id           string  `db:"id"`
	InvoiceId    string  `db:"invoiceId"`
	Description  string  `db:"description"`
	Quantity     int     `db:"quantity"`
	QuantityType string  `db:"quantityType"`
	Price        float64 `db:"price"`
}

type CompanyRecord struct {
	Id                 string `db:"id"`
	Name               string `db:"name"`
	RegistrationNumber string `db:"registrationNumber"`
	CUI                string `db:"CUI"`
	VatID              string `db:"vatID"`
	Email              string `db:"email"`
	Country            string `db:"country"`
	StreetAddress      string `db:"streetAddress"`
	City               string `db:"city"`
	PostCode           string `db:"postCode"`
	BankAccountId      string `db:"bankAccountId"`
}

type ClientRecord struct {
	Id                 string `db:"id"`
	Name               string `db:"name"`
	Country            string `db:"country"`
	RegistrationNumber string `db:"registrationNumber"`
	CUI                string `db:"CUI"`
	VatId              string `db:"vatId"`
	Email              string `db:"email"`
	City               string `db:"city"`
	StreetAddress      string `db:"streetAddress"`
	PostCode           string `db:"postCode"`
	Website            string `db:"website"`
}

type BankAccountRecord struct {
	Id   string `db:"id"`
	Name string `db:"name"`
	IBAN string `db:"iban"`
}

func (r InvoiceRecord) ToDomain() domain.Invoice {
	return domain.Invoice{
		Id:            r.Id,
		InvoiceNumber: r.Id,
		Company:       domain.Company{},
		Client:        domain.Client{},
		Currency:      r.Currency,
		Created:       r.Created.Format(layoutStr),
		Due:           r.Due.Format(layoutStr),
		Articles:      nil,
		Total:         0,
	}
}

func (r BankAccountRecord) ToDomain() domain.BankAccount {
	return domain.BankAccount{
		Name: r.Name,
		IBAN: r.IBAN,
	}
}

func (r CompanyRecord) ToDomain() domain.Company {
	return domain.Company{
		Id:                 r.Id,
		Name:               r.Name,
		RegistrationNumber: r.RegistrationNumber,
		CUI:                r.CUI,
		VatId:              r.VatID,
		Email:              r.Email,
		BankAccount:        domain.BankAccount{},
		BankAccountId:      r.BankAccountId,
		StreetAddress:      r.StreetAddress,
		City:               r.City,
		Country:            r.Country,
		PostCode:           r.PostCode,
	}
}

func (r ArticleRecord) ToDomain() domain.Article {
	var qType domain.QuantityType
	switch r.QuantityType {
	case "unit":
		qType = domain.QuantityTypeUnit
	case "hours":
		qType = domain.QuantityTypeHours
	case "days":
		qType = domain.QuantityTypeDays
	default:
		qType = domain.QuantityTypeUnit
	}
	return domain.Article{
		Description:  r.Description,
		Quantity:     r.Quantity,
		QuantityType: qType,
		Price:        r.Price,
		Amount:       float64(r.Quantity) * r.Price,
	}
}

func (r ArticleRecords) ToDomain() []domain.Article {
	articles := make([]domain.Article, len(r))
	for i, article := range r {
		articles[i] = article.ToDomain()
	}
	return articles
}

func (r ClientRecord) ToDomain() domain.Client {
	return domain.Client{
		Id:                 r.Id,
		Name:               r.Name,
		Country:            r.Country,
		RegistrationNumber: r.RegistrationNumber,
		CUI:                r.CUI,
		VatId:              r.VatId,
		Email:              r.Email,
		City:               r.City,
		PostCode:           r.PostCode,
		Website:            r.Website,
	}
}
