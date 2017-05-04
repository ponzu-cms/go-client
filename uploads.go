package client

import "fmt"

// UploadBySlug makes a GET request to return the Ponzu File Metadata API
// response for the enpoint: `/api/uploads?slug=<Slug>`
func (c *Client) UploadBySlug(slug string) (*APIResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/uploads?slug=%s",
		c.Conf.Host, slug,
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
