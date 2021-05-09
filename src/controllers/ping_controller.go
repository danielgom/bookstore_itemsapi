package controllers

import (
	"github.com/danielgom/bookstore_itemsapi/src/utils/httpUtils"
	"net/http"
)

var HealthController healthControllerInterface = &healthController{}

type healthControllerInterface interface {
	GetHealth(http.ResponseWriter, *http.Request)
}

type healthController struct {
}

func (h *healthController) GetHealth(w http.ResponseWriter, r *http.Request) {
	httpUtils.ResponseJson(w, http.StatusOK, map[string]string{"Status": "Healthy noxus", "request": r.RequestURI})
}
