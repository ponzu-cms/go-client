package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	http.Client
	Conf  *Config
	Cache *Cache
}

type Config struct {
	Host         string `json:"host"`
	DisableCache bool   `json:"disable_cache"`
	Header       http.Header
}

type APIResponse struct {
	Response *http.Response

	JSON []byte
	Data []map[string]interface{}
}

// QueryOptions holds options for a query
type QueryOptions struct {
	Count  int
	Offset int
	Order  string
}

func New(cfg *Config) *Client {
	c := &Client{
		Conf: cfg,
	}

	if cfg.DisableCache {
		return c
	}

	c.Cache = NewCache()
	return c
}

func (c *Client) CacheEnabled() bool {
	return !c.Conf.DisableCache
}

func (a *APIResponse) process() error {
	jsn, err := ioutil.ReadAll(a.Response.Body)
	if err != nil {
		return err
	}

	if len(jsn) == 0 {
		return fmt.Errorf("%s", a.Response.Status)
	}

	data := make(map[string][]map[string]interface{})
	err = json.Unmarshal(jsn, &data)
	if err != nil {
		return err
	}

	a.JSON = jsn
	a.Data = data["data"]

	return nil
}

func setDefaultOpts(opts *QueryOptions) *QueryOptions {
	if opts == nil {
		opts = &QueryOptions{
			Count:  10,
			Order:  "DESC",
			Offset: 0,
		}
	}

	if opts.Count == 0 {
		opts.Count = 10
	}
	if opts.Order == "" {
		opts.Order = "DESC"
	}

	return opts
}

func mergeHeader(req *http.Request, header http.Header) *http.Request {
	for k, v := range header {
		for i := range v {
			req.Header.Add(k, v[i])
		}
	}

	return req
}
