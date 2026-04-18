package api

import (
	"context"
	"net/http"
	"strconv"

	tissue "github.com/mohemohe/go-tissue"
)

type CreateCollectionOption struct {
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private"`
}

type UpdateCollectionOption struct {
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private"`
}

func (c *Client) CreateCollection(ctx context.Context, option *CreateCollectionOption) (*tissue.Collection, error) {
	result := &tissue.Collection{}
	if err := c.sendJSON(ctx, http.MethodPost, "/v1/collections", option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetCollection(ctx context.Context, id int64) (*tissue.Collection, error) {
	result := &tissue.Collection{}
	if err := c.getJSON(ctx, "/v1/collections/"+strconv.FormatInt(id, 10), nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCollection(ctx context.Context, id int64, option *UpdateCollectionOption) (*tissue.Collection, error) {
	result := &tissue.Collection{}
	if err := c.sendJSON(ctx, http.MethodPut, "/v1/collections/"+strconv.FormatInt(id, 10), option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteCollection(ctx context.Context, id int64) error {
	return c.deleteRequest(ctx, "/v1/collections/"+strconv.FormatInt(id, 10))
}
