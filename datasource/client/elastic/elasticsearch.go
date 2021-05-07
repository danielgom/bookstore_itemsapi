package elastic

import (
	"context"
	"encoding/json"
	"github.com/danielgom/bookstore_itemsapi/domain/items"
	"github.com/danielgom/bookstore_utils-go/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
	"time"
)

var EsClient esClientInterface = &esClient{}

type esClientInterface interface {
	setClient(*elasticsearch.Client)
	Index(string, *items.Item) (*esapi.Response, error)
	Get(string, string) (*esapi.Response, error)
}

type esClient struct {
	client *elasticsearch.Client
}

func Init() {

	cfg := elasticsearch.Config{
		Addresses:     []string{"http://localhost:9200"},
		RetryOnStatus: []int{502, 503, 504},
	}

	var err error

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	testConn(client)
	EsClient.setClient(client)
}

func testConn(client *elasticsearch.Client) {

	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting the response from elasticsearch client: %s", err)
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}()

	if res.IsError() {
		log.Fatalf("Error message: %s", res.String())
	}

	logger.Info("Successful connection to elastic search instance")
}

func (c *esClient) setClient(client *elasticsearch.Client) {
	c.client = client
}

func (c *esClient) Index(index string, item *items.Item) (*esapi.Response, error) {

	it, err := json.Marshal(item)
	if err != nil {
		logger.Error("Error trying to parse item", err.Error())
		return nil, err
	}

	res, err := c.client.Index(index, strings.NewReader(string(it)),
		c.client.Index.WithContext(context.Background()),
		c.client.Index.WithDocumentID(item.Id),
		c.client.Index.WithRefresh("true"),
		c.client.Index.WithTimeout(time.Second*1))

	if err != nil {
		logger.Error("Error trying to index document ", err.Error())
		return nil, err
	}

	return res, nil
}

func (c *esClient) Get(index, id string) (*esapi.Response, error) {

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	res, err := c.client.Get(index, id, c.client.Get.WithContext(ctx),
		c.client.Get.WithRefresh(true),
		c.client.Get.WithHuman(),
		c.client.Get.WithErrorTrace())
	if err != nil {
		logger.Error("Error trying to get the document", err.Error())
		return nil, err
	}

	return res, nil
}
