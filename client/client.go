package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const baseURL = "http://127.0.0.1:8080"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type Food struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return errors.Errorf("unknown error, the request failed with %d code", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return errors.New("failed to unmarshall the response")
	}
	fmt.Println(res.StatusCode)
	return nil
}

// Get
func (c *Client) GetFood(ctx context.Context, id int) (*Food, error) {
	q := url.Values{}

	getUrl := fmt.Sprintf("%s/%s/%d", c.BaseURL, "food", id)
	req, err := http.NewRequest("GET", getUrl, strings.NewReader(q.Encode()))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res := Food{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Post
func (c *Client) PostFood(ctx context.Context, food *Food) (*Food, error) {
	u, err := json.Marshal(food)
	if err != nil {
		fmt.Printf("failed to marshal struct %v", err)
		return nil, err
	}

	postUrl := fmt.Sprintf("%s/%s", c.BaseURL, "food")
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(u))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res := Food{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	fmt.Printf("the food id is %d", res.Id)
	return &res, nil
}

// Update
func (c *Client) UpdateFood(ctx context.Context, food *Food, id int) (*Food, error) {
	u, err := json.Marshal(food)
	if err != nil {
		fmt.Printf("failed to marshal struct %v", err)
		return nil, err
	}

	postUrl := fmt.Sprintf("%s/%s/%d", c.BaseURL, "food", id)
	req, err := http.NewRequest("PUT", postUrl, bytes.NewBuffer(u))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res := Food{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Delete
func (c *Client) DeleteFood(ctx context.Context, id int) error {
	var u *Food
	postUrl := fmt.Sprintf("%s/%s/%d", c.BaseURL, "food", id)
	req, err := http.NewRequest("DELETE", postUrl, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	if err := c.sendRequest(req, u); err != nil {
		fmt.Printf("failed to send the DELETE request %v", err)
		return err
	}

	return nil
}
