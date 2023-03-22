package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rsmarincu/billr/builder"
	"github.com/rsmarincu/billr/common"
	"github.com/rsmarincu/billr/domain"
	"github.com/rsmarincu/billr/repository/model"
)

const (
	downloadBillPath = "/downloadBill/{invoiceId}"

	billIdParamName = "invoiceId"
)

type Repository interface {
	GetBill(billId string) (model.InvoiceRecord, error)
}

type Handler interface {
	SetupRouting()
	HandleDownloadBill(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	httpService common.HttpService
	builder     builder.PdfBuilder
}

func NewHandler(service common.HttpService, pdfBuilder builder.PdfBuilder) Handler {
	return &handler{
		builder:     pdfBuilder,
		httpService: service,
	}
}

func (s *handler) HandleDownloadBill(w http.ResponseWriter, r *http.Request) {
	var response domain.Response

	invoiceId, found := getBillIdFromParam(r)
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		response.Status = http.StatusBadRequest
		response.Message = "error"
		response.Data = map[string]interface{}{"data": "bill Id param not set"}
		json.NewEncoder(w).Encode(response)
		return
	}

	pdf, _, err := s.builder.BuildPdf(r.Context(), invoiceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = http.StatusInternalServerError
		response.Message = "error"
		response.Data = map[string]interface{}{"data": fmt.Sprintf("error generating pdf: %v", err)}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(pdf)

	return
}

func (s *handler) SetupRouting() {
	router := s.httpService.GetRouter()

	router.HandleFunc(downloadBillPath, s.HandleDownloadBill).Methods(http.MethodGet)
}

func getBillIdFromParam(r *http.Request) (string, bool) {
	invoiceId := mux.Vars(r)[billIdParamName]
	found := len(invoiceId) > 0

	if found {
		return invoiceId, true
	}

	return "", false
}
