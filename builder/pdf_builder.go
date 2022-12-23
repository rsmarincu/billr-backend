package builder

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/rsmarincu/billr/repository"
)

const gottenburgUrl = "http://localhost:3001/forms/chromium/convert/html"

type PdfBuilder interface {
	BuildPdf(ctx context.Context, billId string) ([]byte, string, error)
}

type pdfBuilder struct {
	repository repository.Repository
	client     http.Client
}

func NewPdfBuilder(repository repository.Repository) PdfBuilder {
	return &pdfBuilder{
		client:     http.Client{},
		repository: repository,
	}
}

func (b *pdfBuilder) BuildPdf(ctx context.Context, invoiceId string) ([]byte, string, error) {
	invoice, err := b.repository.GetInvoice(ctx, invoiceId)
	if err != nil {
		return nil, "", err
	}

	pwd, _ := os.Getwd()
	tmpl := template.Must(template.ParseFiles(pwd + "/internal/invoice_template.html"))

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, invoice)
	if err != nil {
		return nil, "", err
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(fw, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, "", err
	}

	w.Close()
	req, err := http.NewRequest(http.MethodPost, gottenburgUrl, &body)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := b.client.Do(req)
	if err != nil {
		return nil, "", err
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	return respBody, invoice.InvoiceNumber, nil
}
