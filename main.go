package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pienaahj/pbx-api/api"
	"github.com/pienaahj/pbx-api/model"
)

var (
	creds             *model.RequestParams
	user              *model.User
	clientId          string = "JZQulsaolLgZfLLD8SDztJwFn5jlAXWM"
	clientSecrect     string = "ESziOACmxqsHb1SuXcbsBj36rfVdLUU7"
	validToken        string
	validRefreshToken string
)

/*
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
*/

func main() {
	// Create a context
	ctx := context.Background()
	// tr := &http.Transport{
	// 	MaxIdleConns:       10,
	// 	IdleConnTimeout:    30 * time.Second,
	// 	DisableCompression: true,
	// }
	userCreds := &model.UserCreds{
		Username: clientId,
		Password: clientSecrect,
	}

	// client.Client = &http.Client{Transport: tr}
	// client connection established

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	fmt.Println("Connecting to server: ...")
	// fmt.Println("connected to client: ", client)
	// resp, err := client.Get("https://example.com")
	// fmt.Println("client: ", client)

	// get a token
	token, err := api.GetToken(ctx, &client, userCreds)
	if err != nil {
		fmt.Println("Error obtaining token: ", err)
		return
	}
	fmt.Println("access token obtained: ", token)
	fmt.Println()
	// wait for the token
	time.Sleep(5 * time.Second)
	// assign the tokens
	validToken = token.AccessToken
	validRefreshToken = token.RefreshToken
	fmt.Println("valid token: ", validToken)
	fmt.Println("valid refresh token: ", validRefreshToken)
	// get a refresh token
	token, err = api.GetRefreshToken(ctx, &client, validRefreshToken)
	if err != nil {
		fmt.Printf("Error obtaining refresh token: %v", err)
		return
	}

	fmt.Println("refresh token obtained: ", token)

	// assign the tokens
	validToken = token.AccessToken
	validRefreshToken = token.RefreshToken

	//  make a test call
	callReq := &model.CallRequest{
		Caller:         "700",
		Callee:         "0824514478",
		DialPermission: "2002",
		AutoAnswer:     "yes",
	}
	callResp, err := api.MakeCall(ctx, &client, callReq, validToken)
	if err != nil {
		fmt.Printf("Error calling %s: %v\n", callReq.Callee, err)
		return
	}
	fmt.Printf("Calling %v resulted in %s\n", callResp.CallID, callResp.Errmsg)
}
