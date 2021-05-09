package db

import (
	"encoding/json"
	"fmt"
	"github.com/danielgom/bookstore_itemsapi/src/datasource/client/elastic"
	"github.com/danielgom/bookstore_itemsapi/src/domain/items"
	"github.com/danielgom/bookstore_utils-go/errors"
	"github.com/danielgom/bookstore_utils-go/logger"
	"io"
	"strings"
)

type getResponse struct {
	Source items.Item `json:"_source"`
}

const index = "items"

var ItemRepository itemRepositoryInterface = &itemRepository{}

type itemRepositoryInterface interface {
	Save(*items.Item) errors.RestErr
	Get(string) (*items.Item, errors.RestErr)
	SimpleSearch(string) ([]items.Item, error)
}

type itemRepository struct {
}

func (r *itemRepository) Save(i *items.Item) errors.RestErr {
	res, err := elastic.EsClient.Index(index, i)
	if err != nil {
		return errors.NewInternalServerError("error when trying to save the item", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			logger.Error("error closing the response body", err.Error())
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.NewInternalServerError("error trying to read the index response body", err)
	}

	if res.IsError() {
		logger.Error(fmt.Sprintf("error while trying to index document in index %s", i.Id), string(body))
		return errors.NewInternalServerError("error trying to index the document", err)
	}
	return nil
}

func (r *itemRepository) Get(id string) (*items.Item, errors.RestErr) {
	res, err := elastic.EsClient.Get(index, id)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error trying to get the item with id %s", id), err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			logger.Error("error closing the response body", err.Error())
		}
	}()

	var source getResponse

	if err := json.NewDecoder(res.Body).Decode(&source); err != nil {
		logger.Error("invalid json obtained", err.Error())
		return nil, errors.NewInternalServerError("invalid json obtained from the database", err)
	}

	if res.IsError() {
		logger.Error(fmt.Sprintf("error while trying to get the item with id %s", id), res.String())
		if strings.Contains(res.String(), "404") {
			return nil, errors.NewNotFoundError(fmt.Sprintf("document with id %s not found", id))
		}
		return nil, errors.NewInternalServerError("error trying to get the document", err)
	}

	return &source.Source, nil
}

func (r *itemRepository) SimpleSearch(index string) ([]items.Item, error) {

	var query map[string]interface{}

	res, err := elastic.EsClient.Search(index, query)
	if err != nil {
		return nil, errors.NewInternalServerError("error trying to search the items", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			logger.Error("error closing the response body", err.Error())
		}
	}()

	return nil, nil
}
