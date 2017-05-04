package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Content makes a GET request to return the Ponzu Content API response for the
// enpoint: `/api/content?type=<Type>&id=<ID>`
func (c *Client) Content(contentType string, id int) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/content?type=%s&id=%d",
		c.Conf.Host, contentType, id,
	)

	if c.CacheEnabled() {
		ok, resp := c.Cache.Check(endpoint)
		if ok {
			return resp, nil
		}
	}

	w, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.Process()
	if err != nil {
		return resp, err
	}

	if c.CacheEnabled() {
		err := c.Cache.Add(endpoint, resp)
		return resp, err
	}

	return resp, err
}

// Contents makes a GET request to return the Ponzu Content API response for the
// enpoint: `/api/contents?type=<Type>` with query options:
// Count <int>
// Offset <int>
// Order <string>
func (c *Client) Contents(contentType string, opts QueryOptions) (*APIResponse, error) {
	opts = setDefaultOpts(opts)

	endpoint := fmt.Sprintf(
		"%s/api/contents?type=%s&count=%d&offest=%d&order=%s",
		c.Conf.Host, contentType, opts.Count, opts.Offset, opts.Order,
	)

	if c.CacheEnabled() {
		ok, resp := c.Cache.Check(endpoint)
		if ok {
			return resp, nil
		}
	}

	w, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.Process()
	if err != nil {
		return resp, err
	}

	if c.CacheEnabled() {
		err := c.Cache.Add(endpoint, resp)
		return resp, err
	}

	return resp, err
}

func (c *Client) Create(contentType string, data interface{}, fileKeys []string) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/content/create?type=%s",
		c.Conf.Host, contentType,
	)

	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	params := make(url.Values)
	err = json.Unmarshal(j, &params)

	req, err := multipartForm(endpoint, params, fileKeys)
	if err != nil {
		return nil, err
	}

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.Process()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Update(contentType string, id int, data interface{}, fileKeys []string) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/content/update?type=%s&id=%d",
		c.Conf.Host, contentType, id,
	)

	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	params := make(url.Values)
	err = json.Unmarshal(j, &params)

	req, err := multipartForm(endpoint, params, fileKeys)
	if err != nil {
		return nil, err
	}

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.Process()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Delete(contentType string, id int) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/content/delete?type=%s&id=%d",
		c.Conf.Host, contentType, id,
	)

	req, err := multipartForm(endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.Process()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
