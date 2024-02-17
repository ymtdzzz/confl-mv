package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	goconfluence "github.com/virtomize/confluence-go-api"
)

type Client interface {
	GetContentByID(ctx context.Context, id string) (*goconfluence.Content, error)
	GetChildPages(ctx context.Context, id string) (*goconfluence.Search, error)
	MovePage(ctx context.Context, pageID, targetID string) (bool, error)
	CreatePage(ctx context.Context, content *goconfluence.Content) error
}

type client struct {
	api      *goconfluence.API
	urlv1    string
	urlv2    string
	username string
	apiKey   string
}

type count struct {
	id        string
	ancestors []string
	count     int
}

func NewClient(urlv1, urlv2, username, apiKey string) (Client, error) {
	api, err := goconfluence.NewAPI(urlv1, username, apiKey)
	if err != nil {
		return nil, err
	}
	return &client{api, urlv1, urlv2, username, apiKey}, nil
}

func (c *client) GetContentByID(ctx context.Context, id string) (*goconfluence.Content, error) {
	return c.api.GetContentByID(id, goconfluence.ContentQuery{
		Expand: []string{"space", "body.view", "body.storage", "ancestors", "version", "history", "metadata"},
	})
}

func (c *client) GetChildPages(ctx context.Context, id string) (*goconfluence.Search, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/pages/%s/children?limit=250", c.urlv2, id), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.apiKey)
	resp, err := c.api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var search goconfluence.Search
	if err := json.Unmarshal(b, &search); err != nil {
		return nil, err
	}

	return &search, nil
}

func (c *client) MovePage(ctx context.Context, pageID, targetID string) (bool, error) {
	req, _ := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/content/%s/move/append/%s", c.urlv1, pageID, targetID), nil)
	req.SetBasicAuth(c.username, c.apiKey)
	resp, err := c.api.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	return !strings.Contains(string(b), "We can't support moving more than"), nil
}

func (c *client) CreatePage(ctx context.Context, content *goconfluence.Content) error {
	_, err := c.api.CreateContent(content)
	return err
}
