package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/pienaahj/pbx-api/model"
	"golang.org/x/net/websocket"
)

const (
	https_port string = "8088"
	http_port  string = "80"
	//ip address of the PBX server
	pbx_ip       string = "105.246.230.190"
	Api_path     string = "/openapi/v1.0/"
	Content_type string = "application/json"
	ServerAddr   string = "105.246.230.190"
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

func SubscribeToWebsocketService(ctx context.Context, token string) (*websocket.Conn, error) {
	serverUrl := BaseURLSecureWebsocket + Api_path + "subscribe?access_token=" + token
	originUrl := "https://192.168.0.143/"
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
func Handle30008(event []byte) (*model.ExtCallStatus, error) {
	fmt.Println("Handling event30008")
	var (
		extCallStatus              = new(model.ExtCallStatus)
		extentionCallStatusChanged = new(model.ExtentionCallStatusChanged)
	)
	var resp struct {
		Type int    `db:"type" json:"type"`
		SN   string `db:"sn" json:"sn"`
		Msg  string `db:"msg" json:"msg"` // doesn't recognise embeded struct
	}
	err := json.Unmarshal(event, &resp)
	if err != nil {
		log.Printf("Cannot unmarshal from event30015: %v\n", err)
		return extCallStatus, err
	}
	// spew.Dump("Resp: ", resp)
	dataString := resp.Msg
	data := []byte(dataString)
	var messageStruct struct {
		Extension string `db:"extention" json:"extention"`
		Status    string `db:"status" json:"status"`
	}
	err = json.Unmarshal(data, &messageStruct)
	if err != nil {
		log.Printf("Cannot unmarshal from msg: %v\n", err)
		return extCallStatus, err
	}
	// spew.Dump("Msg: ", messageStruct)
	extCallStatus.Extension = messageStruct.Extension
	extCallStatus.Status = messageStruct.Status
	extentionCallStatusChanged.Type = resp.Type
	extentionCallStatusChanged.SN = resp.SN
	extentionCallStatusChanged.Msg = *extCallStatus
	fmt.Printf("Call failed report: %v\n", extCallStatus)
	return extCallStatus, nil
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
func Handle30015(event []byte) (*model.CallFailedReport, error) {
	fmt.Println("Handling event30015")
	var (
		callFailedReport      = new(model.CallFailedReport)
		callFailCallInfo      = new(model.CallFailCallInfo)
		callFailMemebersEmbed = new(model.CallFailMemebers)
		callFailExtention     = new(model.CallFailExtensionInfo)
		callFailInboundInfo   = new(model.CallFailInboundInfo)
		callFailOutboundInfo  = new(model.CallFailOutboundInfo)
	)

	var resp struct {
		Type int    `db:"type" json:"type"`
		SN   string `db:"sn" json:"sn"`
		Msg  string `db:"msg" json:"msg"` // doesn't recognise embeded struct
	}

	// json.RawMessage
	err := json.Unmarshal(event, &resp)
	if err != nil {
		log.Printf("Cannot unmarshal from event30015: %v\n", err)
		return callFailedReport, err
	}
	// spew.Dump("Resp: ", resp)
	dataString := resp.Msg
	data := []byte(dataString)
	var messageStruct struct {
		CallID  string `db:"call_id" json:"call_id"`
		Reason  string `db:"reason," json:"reason"`
		Members string `db:"members,omitempty" json:"members,omitempty"`
	}
	err = json.Unmarshal(data, &messageStruct)
	if err != nil {
		log.Printf("Cannot unmarshal from msg: %v\n", err)
		return callFailedReport, err
	}
	// spew.Dump("Msg: ", messageStruct)
	membersString := messageStruct.Members

	switch membersString {
	case "":
		callFailCallInfo.CallID = messageStruct.CallID
		callFailCallInfo.Reason = messageStruct.Reason
		callFailCallInfo.Members = model.CallFailMemebers{}
		callFailedReport.Type = resp.Type
		callFailedReport.SN = resp.SN
		callFailedReport.Msg = *callFailCallInfo
		spew.Dump(callFailedReport)
		fmt.Println("No memebers found, returning...")
		return callFailedReport, err
	default:
		membersData := []byte(membersString)
		callFailMemebers := struct {
			Extention string `db:"extention" json:"extention"`
			Inbound   string `db:"inbound" json:"inbound"`
			Outbound  string `db:"outbound" json:"outbound"`
		}{}
		err = json.Unmarshal(membersData, &callFailMemebers)
		if err != nil {
			log.Printf("Cannot unmarshal from membersData: %v\n", err)
			return callFailedReport, fmt.Errorf("cannot unmarshal from membersData: %v", err)
		}
		spew.Printf("decoded members: %v of type %T", callFailMemebers)
		// handle Extention
		err = json.Unmarshal([]byte(callFailMemebers.Extention), &callFailExtention)
		if err != nil {
			log.Printf("Cannot unmarshal from Extention: %v\n", err)
			return callFailedReport, fmt.Errorf("cannot unmarshal from Extention: %v", err)
		}
		callFailMemebersEmbed.Extention = *callFailExtention
		// handle Inbound
		err = json.Unmarshal([]byte(callFailMemebers.Inbound), &callFailInboundInfo)
		if err != nil {
			log.Printf("Cannot unmarshal from Inbound: %v\n", err)
			return callFailedReport, fmt.Errorf("cannot unmarshal from Inbound: %v", err)
		}
		callFailMemebersEmbed.Inbound = *callFailInboundInfo
		// handle Outbound
		err = json.Unmarshal([]byte(callFailMemebers.Outbound), &callFailOutboundInfo)
		if err != nil {
			log.Printf("Cannot unmarshal from Outbound: %v\n", err)
			return callFailedReport, fmt.Errorf("cannot unmarshal from Outbound: %v", err)
		}
		callFailMemebersEmbed.Outbound = *callFailOutboundInfo
	}
	// assign the report values
	callFailCallInfo.Members = *callFailMemebersEmbed
	fmt.Printf("Call failed report: %v\n", callFailedReport)
	return callFailedReport, nil
}

// Handle30016(event)
func Handle30016(event []byte) {
	fmt.Println("Handling event30016")
}

// HandleNotImplemented(event)
func HandleNotImplemented(event []byte) {
	fmt.Println("Handling NotImplemented")
}
