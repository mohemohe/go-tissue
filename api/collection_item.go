package api

import (
	"context"
	"net/http"
	"strconv"

	tissue "github.com/mohemohe/go-tissue"
)

type CreateCollectionItemOption struct {
	Link string   `json:"link"`
	Note string   `json:"note,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

type UpdateCollectionItemOption struct {
	Note *string   `json:"note,omitempty"`
	Tags *[]string `json:"tags,omitempty"`
}

func (c *Client) ListCollectionItems(ctx context.Context, collectionID int64, option *PageOption) ([]tissue.CollectionItem, error) {
	query := applyPageOption(nil, option)
	result := []tissue.CollectionItem{}
	path := "/v1/collections/" + strconv.FormatInt(collectionID, 10) + "/items"
	if err := c.getJSON(ctx, path, query, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateCollectionItem(ctx context.Context, collectionID int64, option *CreateCollectionItemOption) (*tissue.CollectionItem, error) {
	result := &tissue.CollectionItem{}
	path := "/v1/collections/" + strconv.FormatInt(collectionID, 10) + "/items"
	if err := c.sendJSON(ctx, http.MethodPost, path, option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateCollectionItem(ctx context.Context, collectionID, itemID int64, option *UpdateCollectionItemOption) (*tissue.CollectionItem, error) {
	result := &tissue.CollectionItem{}
	path := "/v1/collections/" + strconv.FormatInt(collectionID, 10) + "/items/" + strconv.FormatInt(itemID, 10)
	if err := c.sendJSON(ctx, http.MethodPatch, path, option, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteCollectionItem(ctx context.Context, collectionID, itemID int64) error {
	path := "/v1/collections/" + strconv.FormatInt(collectionID, 10) + "/items/" + strconv.FormatInt(itemID, 10)
	return c.deleteRequest(ctx, path)
}
