package app

import (
	"github.com/danielgom/bookstore_itemsapi/datasource/client/elastic"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var router = mux.NewRouter()

func StartApplication() {
	mapUrls()

	elastic.Init()

	srv := &http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
