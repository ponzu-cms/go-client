package client

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type entity struct {
	evictAfter time.Duration
	added      time.Time
	response   *APIResponse
}

type Cache struct {
	mu   sync.Mutex
	Data map[string]entity
}

// func NewCache(evictAfter time.Duration) *Cache {
func NewCache() *Cache {
	c := &Cache{
		mu:   sync.Mutex{},
		Data: make(map[string]entity),
	}

	// TODO: add Data eviction after certain time has passed regardless of
	// entity maxAge attribute

	return c
}

func (c *Cache) Check(endpoint string) (bool, *APIResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	hit, ok := c.Data[endpoint]
	if !ok {
		return false, &APIResponse{}
	}

	if time.Now().After(hit.added.Add(hit.evictAfter)) {
		delete(c.Data, endpoint)
		return false, &APIResponse{}
	}

	return ok, hit.response
}

func (c *Cache) Add(endpoint string, resp *APIResponse) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	maxAge, err := parseMaxAgeHeader(resp)
	if err != nil {
		return err
	}

	c.Data[endpoint] = entity{
		evictAfter: time.Duration(maxAge),
		added:      time.Now(),
		response:   resp,
	}

	fmt.Println(fmt.Sprintf("%#v", c.Data))

	return nil
}

func parseMaxAgeHeader(resp *APIResponse) (time.Duration, error) {
	cc := resp.Response.Header.Get("Cache-Control")
	kv := strings.Split(cc, ", ")
	if len(kv) < 2 {
		return 0, fmt.Errorf("malformed '%s' header", "Cache-Control")
	}

	ma := kv[0]
	parts := strings.Split(ma, "=")
	if len(parts) < 2 {
		return 0, fmt.Errorf("malformed '%s' key", "max-age")
	}

	sec, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Duration(0), err
	}

	return time.Second * time.Duration(sec), nil
}
