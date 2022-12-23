package common

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	EnvKeyServicePort = "HTTP_PORT"
	defaultPort       = "3000"
)

type HttpService interface {
	GetHttpServer() *http.Server
	GetRouter() *mux.Router
	Start()
	Close()
}

type httpService struct {
	name       string
	httpServer *http.Server
	router     *mux.Router
	stopChan   chan error
}

func (h *httpService) GetRouter() *mux.Router {
	return h.router
}

func (h *httpService) GetHttpServer() *http.Server {
	return h.httpServer
}

func (h *httpService) Start() {
	go func() {
		stopper := <-h.stopChan
		log.Printf("stopping service: %v", stopper)
		err := h.httpServer.Close()
		if err != nil {
			log.Println(err)
		}
	}()
}

func (h *httpService) Close() {
	err := h.httpServer.Close()
	if err != nil {
		log.Println(err)
	}
}

func NewHttpService(config Config) HttpService {
	stopChan := make(chan error, 2)

	envPort := config.GetEnvVariable(EnvKeyServicePort, defaultPort)
	port, err := strconv.Atoi(envPort)
	if err != nil {
		panic(err)
	}

	for {
		_, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
		if err != nil {
			break
		} else {
			port += 1
		}
	}

	router := mux.NewRouter()
	handler := cors.AllowAll().Handler(router)
	server := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", port),
	}

	go func() {
		stopChan <- server.ListenAndServe()
	}()

	log.Printf("Internal server listening on: http://127.0.0.1:%d", port)

	return &httpService{
		name:       "HttpService",
		httpServer: server,
		router:     router,
	}
}
