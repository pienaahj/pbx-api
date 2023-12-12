package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt"
	"github.com/pienaahj/pbx-api/api"
	"github.com/pienaahj/pbx-api/model"
	"golang.org/x/net/websocket"
)

var (
	// creds             *model.RequestParams
	// user              *model.User
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

	// make a transport instance
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	// get the client
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	userCreds := &model.UserCreds{
		Username: clientId,
		Password: clientSecrect,
	}
	// MarshalJSONTest()
	fmt.Println("Connecting to server: ...")

	// websocket implementation to monitor events
	// init
	// socketUrl := api.BaseURLSecureWebsocket + api.Api_path + "subscribe?access_token=" + validToken
	// tcpAddr, err := net.ResolveTCPAddr("tcp", api.ServerAddr)
	// if err != nil {
	// 	fmt.Printf("Error resolving websocket server connection: %v\n", api.ServerAddr)
	// }
	// // get a websocket connection
	// conn, err := net.DialTCP("tcp", nil, tcpAddr)
	// if err != nil {
	// 	fmt.Printf("Error dialing websocket server connection: %v\n", api.ServerAddr)
	// }

	// get a token
	token, err := api.GetToken(ctx, client, userCreds)
	if err != nil {
		fmt.Println("Error obtaining token: ", err)
		return
	}
	fmt.Println("access token obtained: ", token)
	fmt.Println()

	// wait for the token
	time.Sleep(2 * time.Second)
	// assign the tokens
	validToken = token.AccessToken
	validRefreshToken = token.RefreshToken
	fmt.Println("valid token: ", validToken)
	// initiate websocket events
	fmt.Println("upgrading network connection to websocket...")
	conn, err := api.SubscribeToWebsocketService(ctx, validToken)
	if err != nil {
		fmt.Printf("Error creating websocket service %v\n", err)
		return
	}
	fmt.Println("successfully subscribed to websocket service")
	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	message := make(chan []byte, 1)
	done := make(chan struct{}, 1)
	go HandleSocketResponse(ctx, done, message, conn)

	// subscribe to the events
	fmt.Println("subscribing events to websocket service...")
	err = api.SubscribeToEvents(ctx, conn, api.Topic_list)
	if err != nil {
		fmt.Println("Error subscribing to events", err)
		done <- struct{}{}
		return
	}
	fmt.Println("successfully subscribed to events")

	//  make a test call
	// callReq := &model.CallRequest{
	// 	Caller:         "700",
	// 	Callee:         "701",
	// 	DialPermission: "700",
	// 	AutoAnswer:     "no",
	// }
	// create a return channel
	returnedCallID := make(chan string, 2)

	// testing make a call
	fmt.Println("****************************************************************")
	fmt.Println("Waiting for incoming events...")
	// service the call id and get the calling number

	// fmt.Println("test call initiated...")
	// callResp, err := api.MakeCall(ctx, client, callReq, validToken)
	// if err != nil {
	// 	fmt.Printf("Error calling %s: %v\n", callReq.Callee, err)
	// 	return
	// }
	// fmt.Printf("Calling %v resulted in %s\n", callResp.CallID, callResp.Errmsg)

	// get the call recording list
	// make a call recording list query
	// callListQuery := &model.QueryRecordingListRequest{
	// 	Page:     1,
	// 	PageSize: 10,
	// 	SortBy:   "id",
	// 	OrderBy:  "asc",
	// }
	// recordingList, err := api.GetRecordingList(ctx, client, callListQuery, validToken)
	// if err != nil {
	// 	fmt.Println("Cannot get recording list: ", err)
	// }
	// // put the call recording list to a file
	// f, err := os.Create("./data/recording_list.txt")
	// if err != nil {
	// 	fmt.Println("Cannot create recording list file", err)
	// }
	// defer f.Close()
	// fmt.Println("recording list message received: ", recordingList.Errmsg)
	// fmt.Println("recording list total pages: ", recordingList.TotalNumber)
	// f.WriteString("Call recording list :\n")
	// for _, recording := range recordingList.Data {
	// 	/*
	// 		ID       int    `db:"id" json:"id"`
	// 		Time     string `db:"time" json:"time"`
	// 		UID      string `db:"uid" json:"uid"`
	// 		CallFrom string `db:"call_from" json:"call_from"`
	// 		CallTo   string `db:"call_to" json:"call_to"`
	// 		Duration int    `db:"duration" json:"duration"`
	// 		Size     int    `db:"size" json:"size"`
	// 		CallType string `db:"call_type" json:"call_type"`
	// 		File     string `db:"file" json:"file"`
	// 	*/
	// 	line := fmt.Sprintf(" ID: %d \t Call from: %s\t Call to: %s\t UID: %s\t Duration: %d Size\t: %d Call type: %s\t File: %s\n", recording.ID, recording.CallFrom, recording.CallTo, recording.UID, recording.Duration, recording.Size, recording.CallType, recording.File)
	// 	f.WriteString(line)
	// }
	// spew.Dump("Recording list response", recordingList)
	// download a recording
	// report := &model.DownloadRecordingRequest{
	// 	ID: 9743,
	// 	// File: "20231002112822-1696238888.81-700-0824514478-Outbound.wav",
	// }
	// err = api.DownloadRecording(ctx, client, report, validToken)
	// if err != nil {
	// 	fmt.Printf("Error downloading recording link: %v\n", err)
	// } else {
	// 	fmt.Println("Successfully downloaded recording wav file")
	// }

	//
	// go func() {
	// 	defer close(done)
	// 	for {
	// 		_, err = conn.Read(msg.Bytes())
	// 		if err != nil {
	// 			fmt.Printf("error reading buffer from websocket service: %v", err)
	// 			return
	// 		}
	// 		fmt.Printf("received message from websocket service: %v\n", msg.String())
	// 	}
	// }()
	for {
		select {
		case <-done:
			return
		case callID := <-returnedCallID: // this now cannot run concurrently
			fmt.Println("received call id on channel")
			// fmt.Println("received the channel for caller id in main:", returnedCallID)
			callerID, err := api.ServiceCallID(ctx, client, callID, validToken)
			if err != nil {
				spew.Dump("error servicing channel call ID: %v", err)
			}
			fmt.Println("received call id: ", callerID)
		case msg := <-message:
			// determine the event type
			var event struct {
				Type int    `json:"type"`
				SN   string `json:"sn"`
				msg  interface{}
			}
			var eventX []byte
			var eventType int
			eventX = msg
			// fmt.Println("received event message from goroutine", string(eventX))
			if eventX != nil {
				err := json.Unmarshal(eventX, &event)
				if err != nil {
					fmt.Printf("error unmarshalling: %v", err)
					eventType = 0
				}
			}
			eventType = event.Type
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
			switch eventType {
			case 30008:
				go api.Handle30008(eventX)
			case 30009:
				go api.Handle30009(eventX)
			case 30011:
				go api.Handle30011(eventX)
			case 30012:
				go api.Handle30012(eventX)
			case 30013:
				go api.Handle30013(eventX)
			case 30014:
				go api.Handle30014(eventX)
			case 30015:
				go api.Handle30015(eventX)
			case 30016:
				// make a caller id channel
				// ch := make(chan string)
				chString, err := api.Handle30016(eventX)
				if err != nil {
					fmt.Println("Failed to handle event 30016")
				}
				fmt.Println("event 30016 handled with callid in handler ", chString)
				fmt.Println("passing caller id to channel")
				// pass the caller id ont the channel to extract the caller
				returnedCallID <- chString

			case 30005:
				api.HandleNotImplemented(eventX)
			case 30006:
				api.HandleNotImplemented(eventX)
			case 30007:
				api.HandleNotImplemented(eventX)
			case 30010:
				api.HandleNotImplemented(eventX)
			case 30017:
				api.HandleNotImplemented(eventX)
			case 30018:
				api.HandleNotImplemented(eventX)
			case 30019:
				api.HandleNotImplemented(eventX)
			case 30020:
				api.HandleNotImplemented(eventX)
			case 30022:
				api.HandleNotImplemented(eventX)
			case 30023:
				api.HandleNotImplemented(eventX)
			case 30024:
				api.HandleNotImplemented(eventX)
			case 0:
				fallthrough
			default:
				fmt.Println("No event type received, waiting...")
			}
		case <-interrupt:
			fmt.Println("Caught interrupt signal - quitting!")
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return
		}
	}

	/*
		fmt.Println("valid refresh token: ", validRefreshToken)
		// get a refresh token
		token, err = api.GetRefreshToken(ctx, client, validRefreshToken)
		if err != nil {
			fmt.Printf("Error obtaining refresh token: %v", err)
			return
		}
		fmt.Println("refresh token obtained: ", token)
		// assign the tokens
		validToken = token.AccessToken
		validRefreshToken = token.RefreshToken
	*/

	// msg := httpsClient("https://localhost:4443")
}

// HandleSocketResponse listens to the responses on the socket and sends it back on the massage channel
func HandleSocketResponse(ctx context.Context, done chan struct{}, message chan []byte, conn *websocket.Conn) {
	defer close(done)
	select {
	case <-done:
		fmt.Println("Done...returning from HandleSocketResponse")
		return
	default:
		for {
			var data []byte
			websocket.Message.Receive(conn, &data)
			if len(data) > 0 {
				spew.Printf("received message from websocket service: %v\n", string(data))
				message <- data
			}
		}
	}
}

// LoadKeys loads the rsa keys from environment variables
func LoadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// load rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := os.ReadFile(privKeyFile)

	if err != nil {
		// fmt.Printf("could not read private key pem file: %v\n", err)
		return nil, nil, fmt.Errorf("could not read private key pem file: %v", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		// fmt.Printf("could not parse private key: %v\n", err)
		return nil, nil, fmt.Errorf("could not parse private key: %v", err)
	}
	fmt.Println("Private key: ", privKey)
	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := os.ReadFile(pubKeyFile)

	if err != nil {
		// fmt.Printf("could not read public key pem file: %v\n", err)
		return privKey, nil, fmt.Errorf("could not read public key pem file: %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		// fmt.Printf("could not parse public key: %v\n", err)
		return privKey, nil, fmt.Errorf("could not parse public key: %v", err)
	}
	fmt.Println("Public Key: ", pubKey)
	return privKey, pubKey, nil
}

// test json marshalling
func MarshalJSONTest() {
	newTest := &model.ExtentionCallStatusChanged{
		Type: 1234,
		SN:   "testSN",
		Msg: model.ExtCallStatus{
			Extension: "test extention",
			Status:    "testing",
		},
	}
	json, err := json.Marshal(newTest)
	if err != nil {
		fmt.Println("Error marshalling test: ", err)
	}
	spew.Dump(json)

}
