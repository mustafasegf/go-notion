package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ketion-so/go-notion/notion/object"
)

const (
	blocksPath = "blocks"
)

// BlocksService handles communication to Notion Blocks API.
//
// API doc: https://developers.notion.com/reference/database
type BlocksService service

// ListBlockChildrenResult object represents the retrieve block children.
//go:generate gomodifytags -file $GOFILE -struct ListBlockChildrenResult -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct ListBlockChildrenResult -add-tags json -w -transform snakecase
type ListBlockChildrenResult struct {
	Object  object.Type   `json:"object"`
	Results []interface{} `json:"results"`
}

type ParagraphBlock struct {
	Object         object.Type
	ID             string
	Type           string
	CreatedTime    string
	LastEditedTime string
	HasChildren    bool
	Text           []RichTextType
	Children       []interface{}
}

type HeadingBlock struct {
	Object         object.Type
	ID             string
	Type           string
	CreatedTime    string
	LastEditedTime string
	HasChildren    bool
	Text           []RichTextType
}

type BulletedListItemBlock struct {
	Object         object.Type
	ID             string
	Type           string
	CreatedTime    string
	LastEditedTime string
	HasChildren    bool
	Text           []RichTextType
	Children       []interface{}
}

type NumberedListItemBlock struct {
	Object         object.Type
	ID             string
	Type           string
	CreatedTime    string
	LastEditedTime string
	HasChildren    bool
	Text           []RichTextType
	Children       []interface{}
}

type TodoBlock struct {
	Object         object.Type
	ID             string
	Type           string
	CreatedTime    string
	LastEditedTime string
	HasChildren    bool
	Text           []RichTextType
	Checked        bool
	Children       []interface{}
}

type ToggleBlock struct {
	Object         object.Type
	ID             string
	Type           string
	CreatedTime    string
	LastEditedTime string
	HasChildren    bool
	Text           []RichTextType
	Children       []interface{}
}

type ChildPageBlock struct {
	Object         object.Type
	ID             string
	Type           string
	CreatedTime    string
	LastEditedTime string
	HasChildren    bool
	Title          string
}

// ListChildren blocks list.
//
// API doc: https://developers.notion.com/reference/get-block-children
func (s *BlocksService) ListChildren(ctx context.Context, blockID string) (*ListBlockChildrenResult, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s/children", blocksPath, blockID))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code not expected, got:%d, message:%s", resp.StatusCode, respErr.Message)
	}

	result := &ListBlockChildrenResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// AppendChildren children block.
//
// API doc: https://developers.notion.com/reference/get-block-children
func (s *BlocksService) AppendChildren(ctx context.Context, blockID string, children interface{}) (interface{}, error) {
	req, err := s.client.NewPostRequest(fmt.Sprintf("%s/%s/children", databasesPath, blockID), children)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code not expected, got:%d, message:%s", resp.StatusCode, respErr.Message)
	}

	var block interface{}
	if err := json.NewDecoder(resp.Body).Decode(&block); err != nil {
		return nil, err
	}

	return block, nil
}
