package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"salesapi/db"
	"salesapi/importer"
	"salesapi/internal/sale"
	appLog "salesapi/log"
)

func main() {
	logger := appLog.InitLogger()

	gormDB := db.Init()

	importerConfig, err := importer.InitConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// data importer service
	importerService := importer.NewImporterService(gormDB, importerConfig, logger)

	// data refresh worker manager to check if worker already running or not
	manager := importer.NewManager()

	// start data refresh worker
	importerService.StartWorker(manager)

	salesAnalytics := sale.NewSaleAnalysisService(gormDB, logger)

	r := mux.NewRouter()

	// logging middleware for all requests
	r.Use(appLog.LoggingMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	api.Handle("/sales/refresh-data", importerService.RefreshHandler(manager)).Methods("POST")
	api.HandleFunc("/sales/top-products", salesAnalytics.GetTopProducts).Methods("GET")

	zap.S().Info("Starting API server at 8080...")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err)
	}
}
