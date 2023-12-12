package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pienaahj/pbx-api/model"
	"golang.org/x/net/websocket"
)

const (
	https_port string = "443"
	http_port  string = "80"
	//ip address of the PBX server
	// pbx_ip       string = "105.246.230.190"
	// pbx_ip string = "13.245.85.196"
	pbx_ip       string = "telectrodev.ras.yeastar.com"
	Api_path     string = "/openapi/v1.0/"
	Content_type string = "application/json"
	// ServerAddr   string = "105.246.230.190"
	ServerAddr string = "telectrodev.ras.yeastar.com"
	// ServerAddr   string = "13.245.85.196"
	Persist      string = "heartbeat"
	SerialNumber string = "3633D2199067"
)

/*
(30008) Extension Call Status Changed	Indicate that the extension call status is changed, and return the current extension call status.
(30009) Extension Presence Status Changed	Indicate that the extension presence status is changed, and return the current extension presence status.
(30011) Call Status Changed 	Indicate that the call status is changed, and return the current call status.
(30012) New CDR	Indicate that a new CDR is generated, and return the call details.
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
		TopicList: []int{30008, 30009, 30011, 30012, 30013, 30014, 30015, 30016},
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
	// spew.Dump(req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error receiving response: ", err)
	}
	defer resp.Body.Close()
	tokenResponse := &model.TokenResponse{}
	// spew.Dump("response", resp.Header)
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
	url := BaseURLSecure + Api_path + "refresh_token"
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

func SubscribeToWebsocketService(ctx context.Context, token string) (*websocket.Conn, error) {
	serverUrl := BaseURLSecureWebsocket + Api_path + "subscribe?access_token=" + token
	originUrl := "https://192.168.0.144/"
	fmt.Println("server url: ", serverUrl)
	originURL, err := url.Parse(originUrl)
	fmt.Println("origin url: ", originURL)
	if err != nil {
		fmt.Printf("Error parsing websocket origin url: %v\n", originUrl)
		return nil, err
	}
	serverURL, err := url.Parse(serverUrl)
	if err != nil {
		fmt.Printf("Error parsing websocket server url: %v\n", serverUrl)
		return nil, err
	}
	// VerifyConnection can be used to replace and customize connection
	// verification. This example shows a VerifyConnection implementation that
	// will be approximately equivalent to what crypto/tls does normally to
	// verify the peer's certificate.

	// Client side configuration.
	tlsConfig := &tls.Config{
		// Set InsecureSkipVerify to skip the default validation we are
		// replacing. This will not disable VerifyConnection.
		InsecureSkipVerify: true,
		// VerifyConnection: func(cs tls.ConnectionState) error {
		// 	opts := x509.VerifyOptions{
		// 		DNSName:       cs.ServerName,
		// 		Intermediates: x509.NewCertPool(),
		// 	}
		// 	for _, cert := range cs.PeerCertificates[1:] {
		// 		opts.Intermediates.AddCert(cert)
		// 	}
		// 	_, err := cs.PeerCertificates[0].Verify(opts)
		// 	return err
		// },
	}

	// Note that when certificates are not handled by the default verifier
	// ConnectionState.VerifiedChains will be nil.
	serverURL.Scheme = "https"
	// construct the websocket config
	myConfig, err := websocket.NewConfig(serverUrl, originUrl)
	if err != nil {
		fmt.Printf("Error creating new config: %v\n", err)
	}
	myConfig.TlsConfig = tlsConfig

	conn, err := websocket.DialConfig(myConfig)
	if err != nil {
		fmt.Printf("Error dialing websocket server connection: %v\n", err)
		return nil, fmt.Errorf("error dialing websocket server connection: %v", err)
	}
	// persist the connection
	_, err = conn.Write([]byte("heartbeat"))
	if err != nil {
		return conn, fmt.Errorf("error persisting socket connection: %v", err)
	}
	var data []byte

	err = websocket.Message.Receive(conn, &data)
	if err != nil {
		return conn, fmt.Errorf("cannot read from buffer: %v", err)
	}
	fmt.Println("Data received: ", string(data))
	fmt.Println("Connection established and persisted successfully")
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
func SubscribeToEvents(ctx context.Context, conn *websocket.Conn, events model.EventTopics) error {
	eventsData, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("cannot marshal events: %v", err)
	}
	_, err = conn.Write(eventsData)
	if err != nil {
		return fmt.Errorf("cannot write events: %v", err)
	}
	var data []byte
	var resp model.SocketSubscriptionResponse
	err = websocket.Message.Receive(conn, &data)
	if err != nil {
		return fmt.Errorf("cannot read from buffer: %v", err)
	}
	fmt.Println("Data received: ", string(data))
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return fmt.Errorf("cannot unmarshal from buffer: %v", err)
	}
	if resp.ErrCode != 0 {
		fmt.Printf("Response received: %v\n", resp)
		return fmt.Errorf("could not subscribe to events: %v", resp.ErrMsg)
	}
	return nil
}

// MakeCall makes a call on the PBX
func MakeCall(ctx context.Context, client *http.Client, call *model.CallRequest, token string) (*model.CallResponse, error) {
	url := BaseURLSecure + Api_path + "call/dial?access_token=" + token
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

// Get call information
func GetCallInfo(ctx context.Context, client *http.Client, callID string, token string) (*model.CallInfoResponse, error) {
	fmt.Println("Getting call info for call id...", callID)
	url := BaseURLSecure + Api_path + "/" + callID + "/recording/list?access_token=" + token

	callRequestData, err := json.Marshal(callID)
	if err != nil {
		fmt.Println("Error marshalling the call id request: ", err)
	}
	// create the credentials string as a Reader
	rdr := bytes.NewBuffer(callRequestData)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, rdr)
	if err != nil {
		fmt.Println("Error creating the call recording report list request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")

	// spew.Dump("Report list request: ", req.Body)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving call recording report list response: ", err)
	}
	defer resp.Body.Close()
	// try to unmarshal it into the response object
	trueResponse := new(model.CallInfoResponse)
	// decode to this struct
	err = json.NewDecoder(resp.Body).Decode(&trueResponse)
	if err != nil {
		fmt.Println("Error decoding call info response: ", err)
	}
	// var callInfo = new(callInfo)
	// trueResp := Sanitize(event)
	// err := json.Unmarshal(trueResp, &callStatusChanged)
	// if err != nil {
	// 	log.Printf("Cannot unmarshal from event30011: %v\n", err)
	// 	return nil, err
	// }
	if trueResponse.ErrCode != 0 {
		fmt.Printf("Error occured requesting call info: %v with message %s\n", err, trueResponse.ErrMsg)
		return &model.CallInfoResponse{}, err
	}
	fmt.Printf("Call info: %v\n", trueResponse)
	return nil, err
}

// serviceCallID processes the call id string returned from events
func ServiceCallID(ctx context.Context, client *http.Client, callID string, token string) (string, error) {
	fmt.Println("ServiceCallID invoked...")

	url := BaseURLSecure + Api_path + "call/query?access_token=" + token
	fmt.Printf("ServiceCallID url: %v\n", url)
	respData := &model.CallQueryResponse{}
	// this should block until the channel has a value
	callData, err := json.Marshal(callID)
	if err != nil {
		fmt.Println("Error marshalling the channel: ", err)
		return "", err
	}
	fmt.Println("channel data", string(callData))
	// create the credentials string as a Reader
	rdr := bytes.NewBuffer(callData)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, rdr)
	if err != nil {
		fmt.Println("Error creating the call id request in serrviceCallID: ", err)
		return "", err
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving query response: ", err)
		return "", err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		fmt.Println("Error decoding json query response body: ", err)
		return "", err
	}
	// check the errcode 0: Success otherwise failed
	if respData.ErrCode != 0 {
		fmt.Printf("Error occured requesting token: %v with message %s\n", err, respData.ErrMsg)
		return "", err
	}
	response := respData.Data[len(respData.Data)-1].Members[len(respData.Data[len(respData.Data)-1].Members)-1].Inbound.From
	fmt.Println("Caller id: ", response)
	return response, nil
}

// Get the recording list
func GetRecordingList(ctx context.Context, client *http.Client, call *model.QueryRecordingListRequest, token string) (*model.QueryRecordingListResponse, error) {
	fmt.Println("Getting call recording list...")
	url := BaseURLSecure + Api_path + "recording/list?access_token=" + token
	// define the response
	// reportResponse := &model.RecordingListResponse{}
	reportData, err := json.Marshal(call)
	if err != nil {
		fmt.Println("Error marshalling the call recording list request: ", err)
	}
	fmt.Println("call recording report list request data", string(reportData))

	// create the credentials string as a Reader
	rdr := bytes.NewBuffer(reportData)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, rdr)
	if err != nil {
		fmt.Println("Error creating the call recording report list request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")

	// spew.Dump("Report list request: ", req.Body)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving call recording report list response: ", err)
	}
	defer resp.Body.Close()

	// try to unmarshal it into the response object
	trueResponse := new(model.QueryRecordingListResponse)
	// decode to this struct
	err = json.NewDecoder(resp.Body).Decode(&trueResponse)
	if err != nil {
		fmt.Println("Error decoding recording list response: ", err)
	}
	/* response format
	ID       int    `db:"id" json:"id"`
	Time     string `db:"time" json:"time"`
	UID      string `db:"uid" json:"uid"`
	CallFrom string `db:"call_from" json:"call_from"`
	CallTo   string `db:"call_to" json:"call_to"`
	Duration int    `db:"duration" json:"duration"`
	Size     int    `db:"size" json:"size"`
	CallType string `db:"call_type" json:"call_type"`
	File     string `db:"file" json:"file"`
	*/

	if trueResponse.Errcode != 0 {
		fmt.Printf("Error occured requesting call recording list: %v with message %s\n", err, trueResponse.Errmsg)
		return &model.QueryRecordingListResponse{}, err
	}
	return trueResponse, nil
}

// Download call recording returns the recording file url
func DownloadRecording(ctx context.Context, client *http.Client, call *model.DownloadRecordingRequest, token string) error {
	fmt.Println("Getting call recording link...")
	url := BaseURLSecure + Api_path + "recording/download?id=9743&access_token=" + token
	// reportResponse := new(model.DownloadRecordingResponse)
	reportResponse := struct {
		Errcode             int    `json:"errcode"`
		Errmsg              string `json:"errmsg"`
		File                string `json:"file"`
		DownloadResourceURL string `json:"download_resource_url"`
	}{}
	// reportData, err := json.Marshal(call)
	// if err != nil {
	// 	fmt.Println("Error marshalling the  report list request: ", err)
	// }
	// fmt.Println("report list request data", string(reportData))
	// create the credentials string as a Reader
	// rdr := bytes.NewBuffer(reportData)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating the report list request: ", err)
	}
	// set the content type on the header
	req.Header.Set("Content-Type", "application/json")

	spew.Dump("Report list request: ", req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error recieving report list response: ", err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&reportResponse)
	if err != nil {
		fmt.Println("Error decoding json response body: ", err)
	}
	// check the errcode 0: Success otherwise failed
	if reportResponse.Errcode != 0 {
		fmt.Printf("Error occured requesting recording link: with message %s\n", reportResponse.Errmsg)
		return err
	}
	// create the wav file
	name, ok := strings.CutPrefix(reportResponse.DownloadResourceURL, "/api/download/")
	if !ok {
		fmt.Println("Error creating wav filename")
		return fmt.Errorf("error creating wav filename from response: %s", reportResponse.DownloadResourceURL)
	}
	filename := fmt.Sprintf("./data/recordings/%s", name)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Cannot create empty wav file: ", err)
	}
	defer f.Close()

	// request the download link
	linkURL := BaseURLSecure + reportResponse.DownloadResourceURL + "?access_token=" + token
	reqLink, err := http.NewRequestWithContext(ctx, http.MethodGet, linkURL, nil)
	fmt.Println("Url link request: ", linkURL)
	if err != nil {
		fmt.Println("Error creating the call recording download request: ", err)
	}
	respLink, err := client.Do(reqLink)
	if err != nil {
		fmt.Println("Error recieving report url response: ", err)
	}
	defer respLink.Body.Close()

	// Put content on file
	body, err := io.ReadAll(respLink.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	// fmt.Printf("body size: %d\n", len(body))
	size, err := f.Write(body)
	// size, err := io.Copy(f, resp.Body) no doesn't work!
	if err != nil {
		fmt.Println("Error decoding json url response body: ", err)
	}
	fmt.Println("Size of wav file downloaded: ", size, " bytes")
	if size <= 0 {
		return fmt.Errorf("error downloading wav file with 0 lenght: %s", reportResponse.DownloadResourceURL)
	}
	return err
}

// event repsonses
// handle event method
// Handle30008(event) extention call status changed
func Handle30008(event []byte) {
	fmt.Println("******************************30008 start *********************")
	fmt.Println("Handling event30008 - extention call status changed")
	var newESC = new(model.ExtentionCallStatusChanged)
	callType, err := newESC.Handle(event)
	if err != nil {
		fmt.Println("Error handling event30008 - extention call status changed", err)
	}
	fmt.Println("Extention status changed to type: ", callType)
	fmt.Println("******************************30008 exit ***********************")
}

// Handle30009(event) extention presence status changed
func Handle30009(event []byte) {
	fmt.Println("******************************30009 start **********************")
	fmt.Println("Handling event30009 - extention presence status changed")
	fmt.Println("******************************30009 exit  ***********************")
}

// Handle30011(event) call status changed
func Handle30011(event []byte) {
	fmt.Println("******************************30011 start **********************")
	fmt.Println("Handling event30011 - call status changed")
	var newCSC = new(model.CallStatusChanged)
	callID, err := newCSC.Handle(event)
	if err != nil {
		fmt.Println("Error handling event30011 - call status changed", err)
	}
	// print the new CDR
	fmt.Println("Call changed call id: ", callID)
	fmt.Println("******************************30011 exit ***********************")
}

// Handle30012(event) new CDR
func Handle30012(event []byte) {
	fmt.Println("******************************30012 start ***********************")
	fmt.Println("Handling event30012 - new CDR")
	var newCDR = new(model.CDR)
	callType, err := newCDR.Handle(event)
	if err != nil {
		fmt.Println("Error handling event30012 - new CDR", err)
	}
	// print the new CDR
	fmt.Println("CDR call type: ", callType)
	fmt.Println("******************************30012 exit ************************")
}

// Handle30013(event) call transfer
func Handle30013(event []byte) {
	fmt.Println("******************************30013 start ***********************")
	fmt.Println("Handling event30013 - call transfer")
	fmt.Println("******************************30013 exit ************************")
}

// Handle30014(event) call forward
func Handle30014(event []byte) {
	fmt.Println("******************************30014 start ***********************")
	fmt.Println("Handling event30014 - call forward")
	fmt.Println("******************************30014 exit  ***********************")
}

// Handle30015(event) Call Failed
func Handle30015(event []byte) {
	fmt.Println("*****************************30015 start ************************")

	var newCFR = new(model.CallFailedReport)
	callID, err := newCFR.Handle(event)
	if err != nil {
		fmt.Println("Error handling event30015", err)
	}

	fmt.Println("Call failed report callid: ", callID)
	fmt.Println("*****************************30015 exit ************************")
}

// Handle30016(event) inbound call invitation
func Handle30016(event []byte) (string, error) {
	fmt.Println("******************************30016 start **********************")
	// ch := make(chan string)
	var chString string = ""
	var newICI = new(model.InboundCallInvitation)
	callID, err := newICI.Handle(event)
	if err != nil {
		fmt.Println("Error handling event 30016", err)
		return chString, err
	}
	// print the new inbound call invitaion
	fmt.Println("Incoming Call Invitation callid returned from api: ", callID)
	chString = callID
	fmt.Println("******************************30016 exit **********************")
	return chString, nil
}

// HandleNotImplemented(event)
func HandleNotImplemented(event []byte) {
	fmt.Println("Handling NotImplemented")
}
