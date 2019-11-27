package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip"
)

// Items ...
type Items struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ToURL ...
func (items *Items) ToURL() string {
	return items.endpoint
}

// Conf ...
func (items *Items) Conf(config *RequestConfig) *Items {
	items.config = config
	return items
}

// Select ...
func (items *Items) Select(oDataSelect string) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$select"] = oDataSelect
	return items
}

// Expand ...
func (items *Items) Expand(oDataExpand string) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$expand"] = oDataExpand
	return items
}

// Filter ...
func (items *Items) Filter(oDataFilter string) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$filter"] = oDataFilter
	return items
}

// Top ...
func (items *Items) Top(oDataTop int) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$top"] = string(oDataTop)
	return items
}

// OrderBy ...
func (items *Items) OrderBy(oDataOrderBy string, ascending bool) *Items {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	if items.modifiers["$orderby"] != "" {
		items.modifiers["$orderby"] += ","
	}
	items.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return items
}

// Get ...
func (items *Items) Get() ([]byte, error) {
	apiURL, _ := url.Parse(items.endpoint)
	query := url.Values{}
	for k, v := range items.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(items.client)
	return sp.Get(apiURL.String(), getConfHeaders(items.config))
}

// Add ...
func (items *Items) Add(body []byte) ([]byte, error) {
	sp := NewHTTPClient(items.client)
	return sp.Post(items.endpoint, body, getConfHeaders(items.config))
}

// GetByID ...
func (items *Items) GetByID(itemID int) *Item {
	return &Item{
		client: items.client,
		config: items.config,
		endpoint: fmt.Sprintf("%s(%d)",
			items.endpoint,
			itemID,
		),
	}
}

// GetByCAML ...
func (items *Items) GetByCAML(caml string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/GetItems", strings.TrimRight(items.endpoint, "/items"))
	apiURL, _ := url.Parse(endpoint)
	query := url.Values{}
	for k, v := range items.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()

	body := trimMultiline(`{
		"query": {
			"__metadata": {"type": "SP.CamlQuery"},
			"ViewXml": "` + trimMultiline(caml) + `"
		}
	}`)

	sp := NewHTTPClient(items.client)
	headers := getConfHeaders(items.config)

	headers["Accept"] = "application/json;odata=verbose"
	headers["Content-Type"] = "application/json;odata=verbose;charset=utf-8"

	return sp.Post(apiURL.String(), []byte(body), headers)
}
