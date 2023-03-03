package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/luizmoitinho/bookstore_oauth/domain"
	"github.com/luizmoitinho/bookstore_oauth/errors"

	"gopkg.in/resty.v1"
)

const (
	TIMEOUT = 100
)

type Methods interface {
	GetAccessToken(string) (*domain.AccessToken, *errors.Rest)
}

type Client struct {
	BaseURL string
	URI     string
	_http   *resty.Client
}

func NewClient(c *resty.Client, baseUrl, uri string) Methods {
	return &Client{
		BaseURL: baseUrl,
		URI:     uri,
		_http:   c,
	}
}

func (c *Client) GetAccessToken(token string) (*domain.AccessToken, *errors.Rest) {
	timeoutDuration := time.Duration(TIMEOUT) * time.Millisecond
	c._http.SetTimeout(timeoutDuration)

	uriWithToken := fmt.Sprintf(c.URI, token)
	resp, err := c._http.R().Get(fmt.Sprintf("%s%s", c.BaseURL, uriWithToken))
	if err != nil {
		return nil, errors.NewInternalServerError("invalid rest client response when trying to get access token")
	}

	if resp.StatusCode() >= http.StatusMultipleChoices {
		var restErr errors.Rest
		if err := json.Unmarshal(resp.Body(), &restErr); err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to get access token")
		}
		return nil, &restErr
	}

	var at domain.AccessToken
	if err := json.Unmarshal(resp.Body(), &at); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unsmarshall access token response")
	}

	return &at, nil
}
