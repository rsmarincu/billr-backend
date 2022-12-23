package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rsmarincu/billr/common"
	"github.com/rsmarincu/billr/domain"
	"github.com/rsmarincu/billr/repository/model"
)

const (
	queryInvoices    = `SELECT "id", "userId", "userCompanyId", "invoicedCompanyId", "currency", "created", "due", "total" FROM "Invoice" WHERE "id"=$1`
	queryCompany     = `SELECT "id", "name", "registrationNumber", "cui", "vatId", "email", "country", "streetAddress", "city", "postCode" FROM "Company" WHERE "id"=$1`
	queryArticles    = `SELECT "id", "invoiceId", "description", "quantity", "quantityType", "price" FROM "Article" WHERE "invoiceId"=$1`
	queryBankAccount = `SELECT "id", "companyId", "name", "IBAN" FROM "BankAccount" WHERE "companyId"=$1`
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

	invoicedCompany, err := r.getCompanyRecord(ctx, invoiceRecord.BilledCompanyId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting company with id: %s: %w", invoiceRecord.BilledCompanyId, err)
	}

	invoicedCompanyBankAccount, err := r.getBankAccount(ctx, invoiceRecord.BilledCompanyId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting bank account: %w", err)
	}

	userCompany, err := r.getCompanyRecord(ctx, invoiceRecord.UserCompanyId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting company with id: %s: %w", invoiceRecord.UserCompanyId, err)
	}

	userCompanyBankAccount, err := r.getBankAccount(ctx, invoiceRecord.UserCompanyId)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting bank account: %w", err)
	}

	articles, err := r.getArticles(ctx, invoiceRecord.Id)
	if err != nil {
		return domain.Invoice{}, fmt.Errorf("error getting articles: %w", err)
	}

	invoice := invoiceRecord.ToDomain()
	invoicedCompany.BankAccount = invoicedCompanyBankAccount
	userCompany.BankAccount = userCompanyBankAccount
	invoice.InvoicedCompany = invoicedCompany
	invoice.UserCompany = userCompany
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
		&invoiceRecord.UserCompanyId,
		&invoiceRecord.BilledCompanyId,
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
		&companyRecord.Email,
		&companyRecord.Country,
		&companyRecord.StreetAddress,
		&companyRecord.City,
		&companyRecord.PostCode,
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

func (r *repository) getBankAccount(ctx context.Context, companyId string) (domain.BankAccount, error) {
	row := r.db.QueryRowContext(ctx, queryBankAccount, companyId)
	var bankAccountRecord model.BankAccountRecord
	if err := row.Scan(
		&bankAccountRecord.Id,
		&bankAccountRecord.CompanyId,
		&bankAccountRecord.Name,
		&bankAccountRecord.IBAN,
	); err != nil {
		return domain.BankAccount{}, err
	}

	return bankAccountRecord.ToDomain(), nil
}
