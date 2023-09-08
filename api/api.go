package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/pienaahj/pbx-api/model"
)

const (
	https_port string = "8088"
	http_port  string = "80"
	//ip address of the PBX server
	pbx_ip       string = "192.168.5.100"
	Api_path     string = "/openapi/v1.0/"
	Content_type string = "application/json"
	ServerAddr   string = "192.168.5.100"
	Persist      string = "heartbeat"
	SerialNumber string = "3633D2199067"
)

/*
(30008) Extension Call Status Changed	Indicate that the extension call status is changed, and return the current extension call status.
(30009) Extension Presence Status Changed	Indicate that the extension presence status is changed, and return the current extension presence status.
(30011) Call Status Changed 	Indicate that the call status is changed, and return the current call status.
(30012) New CDR	Indicate that a new CDR is generated, and return the call de¬ tails.
(30013) Call Transfer	Indicate that a call is transferred, and return the call details.
(30014) Call Forward	Indicate that a call is forwarded, and return the call details.
(30015) Call Failed	Indicate that a call is failed, and return the call details.
(30016) Inbound Call Invitation	Indicate that an inbound call comes from the monitored trunk, and return the call details.
*/

var (
	// for unsecured connections
	BaseURLUnsecure          string = "http://" + pbx_ip + ":" + http_port
	BaseURLUnsecureWebsocket string = "ws://" + pbx_ip + ":" + http_port
	// for secured connections
	BaseURLSecure          string = "https://" + pbx_ip + ":" + https_port
	BaseURLSecureWebsocket string = "wss://" + pbx_ip + ":" + https_port
	connectionURL          string
	access_token           string
	refresh_token          string
	Topic_list             = model.EventTopics{
		TopicList: []string{"3008", "3009", "30011", "30012", "30013", "30014", "30015", "30016"},
	}
)

// The request type we need to make to sign in to the pbx server
type signinReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// The request type we need to sign up onto the pbx sever
type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func GetToken(ctx context.Context, client *http.Client, creds *model.UserCreds) (*model.Tokens, error) {
	// url := BaseURLUnsecure + Api_path + "get_token"
	url := BaseURLSecure + Api_path + "get_token"
	fmt.Println("base url used: ", url)
	tokens := &model.Tokens{}

	credsData, err := json.Marshal(creds)
	if err != nil {
		fmt.Println("Error marshalling the credentials: ", err)
	}
	// fmt.Println("credential data", string(credsData))
	// create the credentials string as a Reader
	rdr := bytes.NewBuffer(credsData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, rdr)
	if err != nil {
		fmt.Println("Error creating the request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")
	// fmt.Printf("request sent: %+v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error receiving response: ", err)
	}
	defer resp.Body.Close()
	tokenResponse := &model.TokenResponse{}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Println("Error decoding json response body: ", err)
	}
	// check the errcode 0: Success otherwise failed
	if tokenResponse.Errcode != 0 {
		fmt.Printf("Error occured requesting token: %v with message %s\n", err, tokenResponse.Errmsg)
		return &model.Tokens{}, err
	}
	fmt.Println("Token expiration time: ", tokenResponse.Access_token_expire_time)
	tokens = &model.Tokens{
		AccessToken:  tokenResponse.Access_token,
		RefreshToken: tokenResponse.Refresh_token,
	}
	return tokens, err
}

func GetRefreshToken(ctx context.Context, client *http.Client, refreshToken string) (*model.Tokens, error) {
	url := BaseURLUnsecure + Api_path + "refresh_token"
	fmt.Println("base url used: ", url)
	tokens := &model.Tokens{}
	// create the refresh token as a string
	refreshStruct := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: refreshToken,
	}
	tokenData, err := json.Marshal(refreshStruct)
	if err != nil {
		fmt.Println("Error marshalling the refresh token: ", err)
	}

	rdr := bytes.NewBuffer(tokenData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, rdr)
	if err != nil {
		fmt.Println("Error creating the request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")
	// fmt.Printf("request sent: %+v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving response: ", err)
	}
	defer resp.Body.Close()
	tokenResponse := &model.TokenResponse{}
	// Errcode:                   resp.Body.errcode,
	// Errmsg:                    resp.Body.errmsg,
	// Access_token_expire_time:  resp.Body.access_token_expire_time,
	// Refresh_token_expire_time: resp.Body.refresh_token_expire_time,
	// Access_token:              resp.Body.access_token,
	// Refresh_token:             resp.Body.refresh_token,
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Println("Error decoding json response body: ", err)
	}
	// check the errcode 0: Success otherwise failed
	if tokenResponse.Errcode != 0 {
		fmt.Printf("Error occured requesting token: %v with message %s\n", err, tokenResponse.Errmsg)
		return &model.Tokens{}, err
	}
	fmt.Println("Token expiration tome: ", tokenResponse.Access_token_expire_time)
	tokens = &model.Tokens{
		AccessToken:  tokenResponse.Access_token,
		RefreshToken: tokenResponse.Refresh_token,
	}
	return tokens, err
}

// SubscribeToWebsocketService subscribe to the websocket service
func SubscribeToWebsocketService(ctx context.Context, token string) (*net.TCPConn, error) {
	url := BaseURLSecureWebsocket + Api_path + "subscribe?access_token=" + token
	fmt.Println("websocket url: ", url)
	// req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	// if err != nil {
	// 	fmt.Println("Error creating the request: ", err)
	// 	return nil, fmt.Errorf("Error creating the request: %v", err)
	// }
	// // set the content type on the header
	// req.Header.Set("Content-Type", "application/json")

	// resp, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println("Error recieving response: ", err)
	// 	return nil, fmt.Errorf("Error recieving responset: %v", err)
	// }
	// defer resp.Body.Close()
	// websocket implementation to monitor events
	// init
	tcpAddr, err := net.ResolveTCPAddr("tcp", url)
	if err != nil {
		fmt.Printf("Error resolving websocket server connection: %v\n", url)
	}
	// get a websocket connection
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("Error dialing websocket server connection: %v\n", url)
		return nil, fmt.Errorf("Error dialing websocket server connection: %v", err)
	}
	// persist the connection
	_, err = conn.Write([]byte("heartbeat"))
	if err != nil {
		return conn, fmt.Errorf("Error pesisting socket connection: %v", err)
	}
	// err = json.NewDecoder(resp.Body).Decode(&callResponse)
	// if err != nil {
	// 	fmt.Println("Error decoding json response body: ", err)
	// }
	// // check the errcode 0: Success otherwise failed
	// if callResponse.Errcode != 0 {
	// 	fmt.Printf("Error occured requesting token: %v with message %s\n", err, callResponse.Errmsg)
	// 	return &model.CallResponse{}, err
	// }

	return conn, nil
}

// SubscribeToEvents subscribe to the websocket service
func SubscribeToEvents(ctx context.Context, conn *net.TCPConn, events model.EventTopics) error {
	eventsString, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("Cannot marshal events: %v", err)
	}
	_, err = conn.Write(eventsString)
	if err != nil {
		return fmt.Errorf("Cannot write events: %v", err)
	}
	var buf []byte
	var resp model.SocketSubscriptionResponse
	_, err = conn.Read(buf)
	if err != nil {
		return fmt.Errorf("Cannot read from buffer: %v", err)
	}
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal from buffer: %v", err)
	}

	return nil
}

// MakeCall makes a call on the PBX
func MakeCall(ctx context.Context, client *http.Client, call *model.CallRequest, token string) (*model.CallResponse, error) {
	url := BaseURLUnsecure + Api_path + "call/dial?access_token=" + token
	callResponse := &model.CallResponse{}

	callData, err := json.Marshal(call)
	if err != nil {
		fmt.Println("Error marshalling the call: ", err)
	}
	fmt.Println("credential data", string(callData))
	// create the credentials string as a Reader
	rdr := bytes.NewBuffer(callData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, rdr)
	if err != nil {
		fmt.Println("Error creating the request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving response: ", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&callResponse)
	if err != nil {
		fmt.Println("Error decoding json response body: ", err)
	}
	// check the errcode 0: Success otherwise failed
	if callResponse.Errcode != 0 {
		fmt.Printf("Error occured requesting token: %v with message %s\n", err, callResponse.Errmsg)
		return &model.CallResponse{}, err
	}

	return callResponse, err
}

// AcceptCall accepts a call on the PBX
func AcceptCall(ctx context.Context, client *http.Client, call *model.CallRequest, token string) (*model.CallResponse, error) {
	url := BaseURLUnsecure + Api_path + "call/dial?access_token=" + token
	callResponse := &model.CallResponse{}

	callData, err := json.Marshal(call)
	if err != nil {
		fmt.Println("Error marshalling the call: ", err)
	}
	fmt.Println("credential data", string(callData))
	// create the credentials string as a Reader
	rdr := bytes.NewBuffer(callData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, rdr)
	if err != nil {
		fmt.Println("Error creating the request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving response: ", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&callResponse)
	if err != nil {
		fmt.Println("Error decoding json response body: ", err)
	}
	// check the errcode 0: Success otherwise failed
	if callResponse.Errcode != 0 {
		fmt.Printf("Error occured requesting token: %v with message %s\n", err, callResponse.Errmsg)
		return &model.CallResponse{}, err
	}

	return callResponse, err
}

// TransferCall accepts a call on the PBX
func TransferCall(ctx context.Context, client *http.Client, call *model.CallRequest, token string) (*model.CallResponse, error) {
	url := BaseURLUnsecure + Api_path + "call/dial?access_token=" + token
	callResponse := &model.CallResponse{}

	callData, err := json.Marshal(call)
	if err != nil {
		fmt.Println("Error marshalling the call: ", err)
	}
	fmt.Println("credential data", string(callData))
	// create the credentials string as a Reader
	rdr := bytes.NewBuffer(callData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, rdr)
	if err != nil {
		fmt.Println("Error creating the request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving response: ", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&callResponse)
	if err != nil {
		fmt.Println("Error decoding json response body: ", err)
	}
	// check the errcode 0: Success otherwise failed
	if callResponse.Errcode != 0 {
		fmt.Printf("Error occured requesting token: %v with message %s\n", err, callResponse.Errmsg)
		return &model.CallResponse{}, err
	}

	return callResponse, err
}

// event repsonses
// Handle30008(event)
func Handle30008(event []byte) {
	fmt.Println("Handling event30008")
}

// Handle30009(event)
func Handle30009(event []byte) {
	fmt.Println("Handling event30009")
}

// Handle30011(event)
func Handle30011(event []byte) {
	fmt.Println("Handling event30011")
}

// Handle30012(event)
func Handle30012(event []byte) {
	fmt.Println("Handling event30012")
}

// Handle30013(event)
func Handle30013(event []byte) {
	fmt.Println("Handling event30013")
}

// Handle30014(event)
func Handle30014(event []byte) {
	fmt.Println("Handling event30014")
}

// Handle30015(event)
func Handle30015(event []byte) {
	fmt.Println("Handling event30015")
}

// Handle30016(event)
func Handle30016(event []byte) {
	fmt.Println("Handling event30016")
}
