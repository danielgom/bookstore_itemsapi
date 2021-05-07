package services

import (
	"github.com/danielgom/bookstore_itemsapi/domain/items"
	"github.com/danielgom/bookstore_itemsapi/repository/db"
	"github.com/danielgom/bookstore_utils-go/errors"
)

var ItemsService itemsServiceInterface = &itemsService{}

type itemsServiceInterface interface {
	Create(*items.Item) (*items.Item, errors.RestErr)
	GetById(string) (*items.Item, errors.RestErr)
}

type itemsService struct {
}

func (i *itemsService) Create(item *items.Item) (*items.Item, errors.RestErr) {
	if err := db.ItemRepository.Save(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (i *itemsService) GetById(id string) (*items.Item, errors.RestErr) {

	item, err := db.ItemRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return item, nil
}
