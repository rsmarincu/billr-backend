package builder

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rsmarincu/billr/domain"
	repository_mock "github.com/rsmarincu/billr/repository/mocks"
)

const layoutStr = "2006-10-02"

func TestPdfBuilder_BuildPdf(t *testing.T) {
	invoiceId := "test"
	invoice := domain.Invoice{
		Id:            invoiceId,
		InvoiceNumber: "testNumber-1234",
		Company: domain.Company{
			Id:                 "company_1",
			Name:               "UserCompany",
			RegistrationNumber: "F32/321/32",
			CUI:                "21312321312",
			VatId:              "RO1230123",
			Email:              "rsmarincu@gmail.com",
			BankAccount: domain.BankAccount{
				Name: "ING",
				IBAN: "RODSA!32323a",
			},
			StreetAddress: "22 scoala de int",
			City:          "Sibiu",
			Country:       "Romania",
			PostCode:      "550005",
		},
		Client: domain.Client{
			Id:                 "company_1",
			Name:               "UserCompany",
			RegistrationNumber: "F32/321/32",
			CUI:                "21312321312",
			VatId:              "RO1230123",
			Email:              "rsmarincu@gmail.com",
			StreetAddress:      "22 scoala de int",
			City:               "Sibiu",
			Country:            "Romania",
			PostCode:           "550005",
		},
		Currency: "EUR",
		Created:  time.Now().Format(layoutStr),
		Due:      time.Now().Format(layoutStr),
		Articles: []domain.Article{
			{
				Description:  "Test article",
				Quantity:     1,
				QuantityType: domain.QuantityTypeUnit,
				Price:        100,
			},
			{
				Description:  "Test article 2",
				Quantity:     1,
				QuantityType: domain.QuantityTypeUnit,
				Price:        100,
			},
		},
		Total: 100,
	}

	ctx := context.Background()
	mockRepository := repository_mock.NewRepositoryMock(t)
	pdfBuilder := NewPdfBuilder(mockRepository)

	mockRepository.On("GetInvoice", ctx, invoiceId).Return(invoice, nil)

	f, _, err := pdfBuilder.BuildPdf(ctx, invoiceId)
	assert.NoError(t, err)

	file, _ := os.Open("test.html")
	defer file.Close()

	file.Write(f)
}
