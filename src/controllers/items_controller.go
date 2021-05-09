package controllers

import (
	"encoding/json"
	"github.com/danielgom/bookstore_itemsapi/src/domain/items"
	"github.com/danielgom/bookstore_itemsapi/src/services"
	"github.com/danielgom/bookstore_itemsapi/src/utils/httpUtils"
	"github.com/danielgom/bookstore_oauth-go/oauth"
	"github.com/danielgom/bookstore_utils-go/errors"
	"github.com/danielgom/bookstore_utils-go/logger"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

var ItemsController itemsControllerInterface = &itemsController{}

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}

type itemsController struct {
}

func (c *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		httpUtils.ResponseError(w, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		restErr := errors.NewBadRequestError("Invalid request body")
		httpUtils.ResponseError(w, restErr)
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal()
		}
	}()

	item := new(items.Item)
	if err := json.Unmarshal(body, item); err != nil {
		logger.Error("Invalid json body", err.Error())
		restErr := errors.NewBadRequestError("Invalid json body")
		httpUtils.ResponseError(w, restErr)
		return
	}

	item.Seller = oauth.GetCallerId(r)

	result, createErr := services.ItemsService.Create(item)
	if createErr != nil {
		httpUtils.ResponseError(w, createErr)
		return
	}

	httpUtils.ResponseJson(w, http.StatusCreated, result)
}

func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iId := strings.TrimSpace(vars["id"])

	item, getErr := services.ItemsService.GetById(iId)
	if getErr != nil {
		httpUtils.ResponseError(w, getErr)
		return
	}
	httpUtils.ResponseJson(w, http.StatusOK, item)
}
