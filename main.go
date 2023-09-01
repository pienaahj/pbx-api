package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pienaahj/pbx-api/api"
	"github.com/pienaahj/pbx-api/model"
)

type API_client struct {
	Client *http.Client
}

var client *API_client
var creds *model.RequestParams
var user *model.User
var username string = "700"
var password string = "1uuITucqDNbSz7S"

func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	user = &model.User{
		Username: username,
		Password: password,
	}
	client.Client = &http.Client{Transport: tr}
	// resp, err := client.Get("https://example.com")

}

func (c *API_client) GetToken(ctx context.Context, creds *model.RequestParams) (*model.Tokens, error) {
	url := api.BaseURLUnsecure + api.Api_path + "get_token"
	JSONBody, err := json.Marshal(creds)
	if err != nil {
		fmt.Println("Error marshalling credentials: ", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, JSONBody)
	resp, err := c.Client.Do(req)
}

func (c *API_client) GetRefreshToken(token string) (*model.Tokens, error) {}
