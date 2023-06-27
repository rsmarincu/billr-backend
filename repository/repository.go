package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/rsmarincu/billr/common"
	"github.com/rsmarincu/billr/domain"
	"github.com/rsmarincu/billr/repository/model"
)

const (
	queryInvoices    = `SELECT "id", "userId", "companyId", "clientId", "currency", "created", "due", "total" FROM "Invoice" WHERE "id"=$1`
	queryCompany     = `SELECT "id", "name", "registrationNumber", "taxId", "vatId", "country", "streetAddress", "city", "postCode", "bankAccountId" FROM "Company" WHERE "id"=$1`
	queryArticles    = `SELECT "id", "invoiceId", "description", "quantity", "quantityType", "price" FROM "Article" WHERE "invoiceId"=$1`
	queryBankAccount = `SELECT "id",  "name", "iban", "swift" FROM "BankAccount" WHERE "id"=$1`
	queryClient      = `SELECT "id", "name", "country","registrationNumber", "taxId", "vatId", COALESCE(email, '') as email, "streetAddress", "city",  "postCode", "website" FROM "Client" WHERE "id"=$1`
)

type Repository interface {
	GetInvoice(ctx context.Context, invoiceId string) (domain.Invoice, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(service common.DatabaseService) Repository {
	return &repository{
		db: service.GetDb(),
	}
}

func (r *repository) GetInvoice(ctx context.Context, invoiceId string) (domain.Invoice, error) {
	invoiceRecord, err := r.getInvoiceRecord(ctx, invoiceId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting invoice: %w", err)
	}

	client, err := r.getClientRecord(ctx, invoiceRecord.ClientId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting client with id: %s: %w", invoiceRecord.ClientId, err)
	}

	userCompany, err := r.getCompanyRecord(ctx, invoiceRecord.CompanyId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting company with id: %s: %w", invoiceRecord.CompanyId, err)
	}
	log.Println(userCompany)
	userCompanyBankAccount, err := r.getBankAccount(ctx, userCompany.BankAccountId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting bank account: %w", err)
	}

	articles, err := r.getArticles(ctx, invoiceRecord.Id)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting articles: %w", err)
	}

	invoice := invoiceRecord.ToDomain()
	userCompany.BankAccount = userCompanyBankAccount
	invoice.Client = client
	invoice.Company = userCompany
	invoice.Articles = articles

	var total float64
	for _, article := range invoice.Articles {
		total += article.Amount
	}
	invoice.Total = total

	return invoice, nil
}

func (r *repository) getInvoiceRecord(ctx context.Context, invoiceId string) (model.InvoiceRecord, error) {
	row := r.db.QueryRowContext(ctx, queryInvoices, invoiceId)
	var invoiceRecord model.InvoiceRecord
	if err := row.Scan(
		&invoiceRecord.Id,
		&invoiceRecord.UserId,
		&invoiceRecord.CompanyId,
		&invoiceRecord.ClientId,
		&invoiceRecord.Currency,
		&invoiceRecord.Created,
		&invoiceRecord.Due,
		&invoiceRecord.Total,
	); err != nil {
		return model.InvoiceRecord{}, err
	}

	return invoiceRecord, nil
}

func (r *repository) getCompanyRecord(ctx context.Context, companyId string) (domain.Company, error) {
	row := r.db.QueryRowContext(ctx, queryCompany, companyId)
	var companyRecord model.CompanyRecord
	if err := row.Scan(
		&companyRecord.Id,
		&companyRecord.Name,
		&companyRecord.RegistrationNumber,
		&companyRecord.CUI,
		&companyRecord.VatID,
		&companyRecord.Country,
		&companyRecord.StreetAddress,
		&companyRecord.City,
		&companyRecord.PostCode,
		&companyRecord.BankAccountId,
	); err != nil {
		return domain.Company{}, err
	}

	return companyRecord.ToDomain(), nil
}

func (r *repository) getArticles(ctx context.Context, invoiceId string) ([]domain.Article, error) {
	rows, err := r.db.QueryContext(ctx, queryArticles, invoiceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articleRecords model.ArticleRecords

	for rows.Next() {
		var articleRecord model.ArticleRecord
		if err := rows.Scan(
			&articleRecord.Id,
			&articleRecord.InvoiceId,
			&articleRecord.Description,
			&articleRecord.Quantity,
			&articleRecord.QuantityType,
			&articleRecord.Price,
		); err != nil {
			return nil, err
		}
		articleRecords = append(articleRecords, articleRecord)
	}

	return articleRecords.ToDomain(), nil
}

func (r *repository) getBankAccount(ctx context.Context, bankAccountId string) (domain.BankAccount, error) {
	row := r.db.QueryRowContext(ctx, queryBankAccount, bankAccountId)
	var bankAccountRecord model.BankAccountRecord
	if err := row.Scan(
		&bankAccountRecord.Id,
		&bankAccountRecord.Name,
		&bankAccountRecord.IBAN,
		&bankAccountRecord.Swift,
	); err != nil {
		return domain.BankAccount{}, err
	}

	return bankAccountRecord.ToDomain(), nil
}

func (r *repository) getClientRecord(ctx context.Context, companyId string) (domain.Client, error) {
	row := r.db.QueryRowContext(ctx, queryClient, companyId)
	var clientRecord model.ClientRecord
	if err := row.Scan(
		&clientRecord.Id,
		&clientRecord.Name,
		&clientRecord.Country,
		&clientRecord.RegistrationNumber,
		&clientRecord.CUI,
		&clientRecord.VatId,
		&clientRecord.Email,
		&clientRecord.StreetAddress,
		&clientRecord.City,
		&clientRecord.PostCode,
		&clientRecord.Website,
	); err != nil {
		return domain.Client{}, err
	}

	return clientRecord.ToDomain(), nil
}
