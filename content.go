package client

import (
	"fmt"
	"net/http"
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

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req = mergeHeader(req, c.Conf.Header)

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.process()
	if err != nil {
		return resp, err
	}

	if c.CacheEnabled() {
		err := c.Cache.Add(endpoint, resp)
		return resp, err
	}

	return resp, nil
}

// ContentBySlug makes a GET request to return the Ponzu Content API response for the
// enpoint: `/api/content?slug=<Slug>`
func (c *Client) ContentBySlug(slug string) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/content?slug=%s",
		c.Conf.Host, slug,
	)

	if c.CacheEnabled() {
		ok, resp := c.Cache.Check(endpoint)
		if ok {
			return resp, nil
		}
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req = mergeHeader(req, c.Conf.Header)

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.process()
	if err != nil {
		return resp, err
	}

	if c.CacheEnabled() {
		err := c.Cache.Add(endpoint, resp)
		return resp, err
	}

	return resp, nil
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

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req = mergeHeader(req, c.Conf.Header)

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.process()
	if err != nil {
		return resp, err
	}

	if c.CacheEnabled() {
		err := c.Cache.Add(endpoint, resp)
		return resp, err
	}

	return resp, nil
}

// Create makes a POST request containing a multipart/form-data body with the
// contents of a content item to be created and stored in Ponzu. Note: the fileKeys
// []string argument should contain the field/key names of the item's uploads.
//
// The *APIResponse will indicate whether the request failed or succeeded based
// on the contents of its Data or JSON fields, or by checking the Status of it's
// original http.Response. Callers should expect failures to occur when a Content
// type does not implement the api.Createable interface
func (c *Client) Create(contentType string, data url.Values, fileKeys []string) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/content/create?type=%s",
		c.Conf.Host, contentType,
	)

	req, err := multipartForm(endpoint, data, fileKeys)
	if err != nil {
		return nil, err
	}

	req = mergeHeader(req, c.Conf.Header)

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.process()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update makes a POST request containing a multipart/form-data body with the
// contents of a content item to be updated and stored in Ponzu. Note: the fileKeys
// []string argument should contain the field/key names of the item's uploads
//
// The *APIResponse will indicate whether the request failed or succeeded based
// on the contents of its Data or JSON fields, or by checking the Status of it's
// original http.Response. Callers should expect failures to occur when a Content
// type does not implement the api.Updateable interface
func (c *Client) Update(contentType string, id int, data url.Values, fileKeys []string) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/content/update?type=%s&id=%d",
		c.Conf.Host, contentType, id,
	)

	req, err := multipartForm(endpoint, data, fileKeys)
	if err != nil {
		return nil, err
	}

	req = mergeHeader(req, c.Conf.Header)

	w, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &APIResponse{
		Response: w,
	}

	err = resp.process()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete makes a POST request to the proper endpoint and with the required data
// to remove content from Ponzu
//
// The *APIResponse will indicate whether the request failed or succeeded based
// on the contents of its Data or JSON fields, or by checking the Status of it's
// original http.Response. Callers should expect failures to occur when a Content
// type does not implement the api.Deleteable interface
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

	req = mergeHeader(req, c.Conf.Header)

	resp := &APIResponse{
		Response: w,
	}

	err = resp.process()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Reference is a helper method to fetch a content item that is referenced from
// a parent content type
func (c *Client) Reference(uri string) (*APIResponse, error) {
	target, err := ParseReferenceURI(uri)
	if err != nil {
		return nil, err
	}

	return c.Content(target.Type, target.ID)
}
