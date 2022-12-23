package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rsmarincu/billr/api"
	"github.com/rsmarincu/billr/builder"
	"github.com/rsmarincu/billr/common"
	"github.com/rsmarincu/billr/repository"
)

const (
	serviceName = "BillrService"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config := common.NewConfig(serviceName)
	httpService := common.NewHttpService(config)
	databaseService := common.NewDatabaseService(config)

	repo := repository.NewRepository(databaseService)
	pdfBuilder := builder.NewPdfBuilder(repo)
	handler := api.NewHandler(httpService, pdfBuilder)
	billrService := NewBillrService(handler)

	billrService.Init()

	<-done
	log.Println("server stopped")
}
