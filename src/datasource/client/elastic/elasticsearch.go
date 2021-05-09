package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/danielgom/bookstore_itemsapi/src/domain/items"
	"github.com/danielgom/bookstore_utils-go/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"os"
	"strings"
	"time"
)

const (
	responseError      = "error getting the response"
	elasticsearchHost  = "ELASTIC_SEARCH_HOST"
	elasticsearchPorts = "ELASTIC_SEARCH_PORTS"
)

var (
	EsClient esClientInterface = &esClient{}
	host                       = os.Getenv(elasticsearchHost)
	ports                      = os.Getenv(elasticsearchPorts)
)

type esClientInterface interface {
	setClient(*elasticsearch.Client)
	Index(string, *items.Item) (*esapi.Response, error)
	Get(string, string) (*esapi.Response, error)
	Search(string, map[string]interface{}) (*esapi.Response, error)
}

type esClient struct {
	client *elasticsearch.Client
}

func Init() {

	cfg := elasticsearch.Config{
		Addresses:     []string{fmt.Sprintf("http://%s:%s", host, ports)},
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
		logger.Error(responseError, err.Error())
		return nil, err
	}

	return res, nil
}

func (c *esClient) Get(index, id string) (*esapi.Response, error) {

	// Create context with timeout, context should come from the function that is calling the es client methods
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	res, err := c.client.Get(index, id, c.client.Get.WithContext(ctx),
		c.client.Get.WithRefresh(true),
		c.client.Get.WithHuman(),
		c.client.Get.WithErrorTrace())
	if err != nil {
		logger.Error(responseError, err.Error())
		return nil, err
	}

	return res, nil
}

func (c *esClient) Search(index string, query map[string]interface{}) (*esapi.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		logger.Error("error encoding query %s", err.Error())
		return nil, err
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(index),
		c.client.Search.WithBody(&buf),
		c.client.Search.WithTrackTotalHits(true),
		c.client.Search.WithPretty())

	if err != nil {
		logger.Error(responseError, err.Error())
		return nil, err
	}
	return res, nil
}
