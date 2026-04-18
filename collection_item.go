package go_tissue

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type ListCollectionItemsOption struct {
	CollectionID int64
	Page         int
	PerPage      int
}

type CreateCollectionItemOption struct {
	CollectionID int64    `json:"-"`
	Link         string   `json:"link"`
	Note         string   `json:"note,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

type UpdateCollectionItemOption struct {
	CollectionID int64     `json:"-"`
	ItemID       int64     `json:"-"`
	Note         *string   `json:"note,omitempty"`
	Tags         *[]string `json:"tags,omitempty"`
}

func (c *Client) ListCollectionItems(ctx context.Context, option *ListCollectionItemsOption) ([]CollectionItem, error) {
	if option == nil {
		return nil, errNilOption
	}
	query := url.Values{}
	if option.Page > 0 {
		query.Set("page", strconv.Itoa(option.Page))
	}
	if option.PerPage > 0 {
		query.Set("per_page", strconv.Itoa(option.PerPage))
	}
	result := []CollectionItem{}
	path := "/api/collections/" + strconv.FormatInt(option.CollectionID, 10) + "/items"
	if err := c.getJSON(ctx, path, query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateCollectionItem(ctx context.Context, option *CreateCollectionItemOption) (*CollectionItem, error) {
	if option == nil {
		return nil, errNilOption
	}
	result := &CollectionItem{}
	path := "/api/collections/" + strconv.FormatInt(option.CollectionID, 10) + "/items"
	if err := c.sendJSON(ctx, http.MethodPost, path, option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCollectionItem(ctx context.Context, option *UpdateCollectionItemOption) (*CollectionItem, error) {
	if option == nil {
		return nil, errNilOption
	}
	result := &CollectionItem{}
	path := "/api/collections/" + strconv.FormatInt(option.CollectionID, 10) + "/items/" + strconv.FormatInt(option.ItemID, 10)
	if err := c.sendJSON(ctx, http.MethodPatch, path, option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteCollectionItem(ctx context.Context, collectionID, itemID int64) error {
	path := "/api/collections/" + strconv.FormatInt(collectionID, 10) + "/items/" + strconv.FormatInt(itemID, 10)
	return c.deleteRequest(ctx, path)
}
