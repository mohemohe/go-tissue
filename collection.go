package go_tissue

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type CreateCollectionOption struct {
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private"`
}

type UpdateCollectionOption struct {
	ID        int64  `json:"-"`
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private"`
}

type ListCollectionsOption struct {
	Page    int
	PerPage int
}

func (c *Client) ListCollections(ctx context.Context, option *ListCollectionsOption) ([]Collection, error) {
	query := url.Values{}
	if option != nil {
		if option.Page > 0 {
			query.Set("page", strconv.Itoa(option.Page))
		}
		if option.PerPage > 0 {
			query.Set("per_page", strconv.Itoa(option.PerPage))
		}
	}
	result := []Collection{}
	if err := c.getJSON(ctx, "/api/collections", query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateCollection(ctx context.Context, option *CreateCollectionOption) (*Collection, error) {
	result := &Collection{}
	if err := c.sendJSON(ctx, http.MethodPost, "/api/collections", option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCollection(ctx context.Context, option *UpdateCollectionOption) (*Collection, error) {
	if option == nil {
		return nil, errNilOption
	}
	result := &Collection{}
	path := "/api/collections/" + strconv.FormatInt(option.ID, 10)
	if err := c.sendJSON(ctx, http.MethodPut, path, option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteCollection(ctx context.Context, collectionID int64) error {
	return c.deleteRequest(ctx, "/api/collections/"+strconv.FormatInt(collectionID, 10))
}
