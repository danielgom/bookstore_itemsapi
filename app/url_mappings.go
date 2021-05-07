package app

import (
	"github.com/danielgom/bookstore_itemsapi/controllers"
	"net/http"
)

func mapUrls() {
	router.HandleFunc("/items", controllers.ItemsController.Create).Methods(http.MethodPost)
	router.HandleFunc("/items/{id}", controllers.ItemsController.Get).Methods(http.MethodGet)
	router.HandleFunc("/health", controllers.HealthController.GetHealth).Methods(http.MethodGet)
}
