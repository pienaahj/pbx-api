package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
)

// event topics
type Event struct {
	Event30008 ExtCallStatus
	Event30009 ExtensionPresenceStatusChanged
	Event30011 TrunkRegistrationStatusChanged
	Event30012 CallStatusChanged
	Event30014 CDR
	Event30015 CallForwardReport
	Event30016 CallFailedReport
}
type EventHandler interface {
	Handle()
}

// event subscription
//

type CallType struct {
	Extension string `db:"extension" json:"extension"`
}

type CallRequest struct {
	Caller         string `db:"caller" json:"caller"`
	Callee         string `db:"callee" json:"callee"`
	DialPermission string `db:"dial_permission" json:"dial_permission"`
	AutoAnswer     string `db:"auto_answer" json:"auto_answer"`
}
type CallResponse struct {
	Errcode int    `db:"errcode" json:"errcode"`
	Errmsg  string `db:"errmsg" json:"errmsg"`
	CallID  string `db:"call_id" json:"call_id"`
}

type CallAcceptRequest struct {
	ChannelId string `db:"channel_id" json:"channel_id"`
}

type CallAcceptResponse struct {
	Errcode int    `db:"errcode" json:"errcode"`
	Errmsg  string `db:"errmsg" json:"errmsg"`
	CallID  string `db:"call_id" json:"call_id"`
}
type EventTopics struct {
	TopicList []int `db:"topic_list" json:"topic_list"`
}

// event subscription response
type SocketSubscriptionResponse struct {
	ErrCode int    `db:"errcode" json:"errcode"`
	ErrMsg  string `db:"errmsg" json:"errmsg"`
}

// (30008) Extension Call Status Changed
type ExtCallStatus struct {
	Extension string `db:"extension" json:"extension"`
	Status    string `db:"status" json:"status"`
}

type ExtentionCallStatusChanged struct {
	Type int           `db:"type" json:"type"`
	SN   string        `db:"sn" json:"sn"`
	Msg  ExtCallStatus `db:"msg" json:"msg"`
}

// Handle30008(event) extention call status changed
func (e *ExtentionCallStatusChanged) Handle(event []byte) {
	fmt.Println("Handling event30008")
	//{"type":30015,"sn":"3633D2199067","msg":"{\"call_id\":\"1694606403.000055\",\"reason\":\"NOT found\",\"members\":null}"}
	var extentionCallStatusChanged = new(ExtentionCallStatusChanged)
	trueResp := Sanitize(event)
	err := json.Unmarshal(trueResp, &extentionCallStatusChanged)
	if err != nil {
		log.Printf("Cannot unmarshal from event3008: %v\n", err)
		return
	}
	// spew.Dump("extentionCallStatusChanged: ", extentionCallStatusChanged)

	fmt.Printf("Call extention status change report: %v\n", extentionCallStatusChanged)
}

// (30009) Extension Presence Status Changed
type ExtentionPresence struct {
	Extension string `db:"extension" json:"extension"`
	Status    string `db:"status" json:"status"`
}
type ExtensionPresenceStatusChanged struct {
	Type int               `db:"type" json:"type"`
	SN   string            `db:"sn" json:"sn"`
	Msg  ExtentionPresence `db:"msg" json:"msg"`
}

// (30010) Trunk Registration Status Changed
type TrunkRegInfo struct {
	TrunkName    string `db:"trunc_name" json:"trunc_name"`
	Kind         string `db:"kind" json:"kind"`
	Status       int    `db:"status" json:"status"`
	RegisteredIP string `db:"registered_ip" json:"registered_ip"`
}
type TrunkRegistrationStatusChanged struct {
	Type int          `db:"type" json:"type"`
	SN   string       `db:"sn" json:"sn"`
	Msg  TrunkRegInfo `db:"msg" json:"msg"`
}

var TrunkStatus = map[int]string{
	0:  "Unknown status.",
	1:  "The trunk is Idle.",
	2:  "The trunk is Busy.",
	3:  "The SIP trunk is idle and unmonitored.",
	4:  "The SIP trunk is registering.",
	41: "SIP register trunk registration failed.",
	42: "The SIP trunk is unreachable.",
	43: "The SIP account trunk is unavailable.",
	44: "The SIP trunk is disabled.",
}

// (30011) Call Status Changed
type CallStatusOutboundInfo struct {
	From         string `db:"from" json:"from"`
	To           string `db:"to" json:"to"`
	TrunkName    string `db:"trunk_name" json:"trunk_name"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallStatusInboundInfo struct {
	From         string `db:"from" json:"from"`
	To           string `db:"to" json:"to"`
	TrunkName    string `db:"trunk_name" json:"trunk_name"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallStatusExtensionInfo struct {
	Number       string `db:"number" json:"number"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallStatusMembers struct {
	Extention CallStatusExtensionInfo `db:"extension" json:"extension"`
	Inbound   CallStatusInboundInfo   `db:"inbound" json:"inbound"`
	Outbound  CallStatusOutboundInfo  `db:"outbound" json:"outbound"`
}
type CallStatusCallInfo struct {
	CallID  string              `db:"call_id" json:"call_id"`
	Members []CallStatusMembers `db:"members" json:"members"`
}
type CallStatusChanged struct {
	Type int                `db:"type" json:"type"`
	SN   string             `db:"sn" json:"sn"`
	Msg  CallStatusCallInfo `db:"msg" json:"msg"`
}

// Handle30011(event) call status changed
func (e *CallStatusChanged) Handle(event []byte) {
	fmt.Println("Handling event30011")
	/* "type":30011,
	"sn":"3633D2199067",
	"msg":"{\"call_id\":\"1694612246.86\",
	\"members\":
	[{\"extension\":{\"number\":\"700\",
		\"channel_id\":\"PJSIP/700-00000023\",
		\"member_status\":\"BYE\",
		\"call_path\":\"\"}}
	]}
	"}
	*/
	var callStatusChanged = new(CallStatusChanged)
	trueResp := Sanitize(event)
	err := json.Unmarshal(trueResp, &callStatusChanged)
	if err != nil {
		log.Printf("Cannot unmarshal from event30011: %v\n", err)
		return
	}

	fmt.Printf("Call extention status change report: %v\n", callStatusChanged)
}

// (30012) New CDR
type CDRCallDetails struct {
	CallID        string `db:"call_id" json:"call_id"`
	TimeStart     string `db:"time_start" json:"time_start"`
	CallFrom      string `db:"call_from" json:"call_from"`
	CallTo        string `db:"call_to" json:"call_to"`
	CallDuration  string `db:"call_duration" json:"call_duration"`
	TalkDuration  string `db:"talk_duration" json:"talk_duration"`
	SRCTrunkname  string `db:"src_trunk_name" json:"src_trunk_name"`
	DSTTrunkName  string `db:"dst_trunk_name" json:"dst_trunk_name"`
	PinCode       string `db:"pin_code" json:"pin_code"`
	Status        string `db:"status" json:"status"`
	Type          string `db:"type" json:"type"`
	Recording     string `db:"recording" json:"recording"`
	DIDNumber     string `db:"did_number" json:"did_number"`
	AgentRingTime int    `db:"agent_ring_time" json:"agent_ring_time"`
}
type CDR struct {
	Type int            `db:"type" json:"type"`
	SN   string         `db:"sn" json:"sn"`
	Msg  CDRCallDetails `db:"msg" json:"msg"`
}

// Handle30012(event) New CDR Event
func (c *CDR) Handle(event []byte) error {
	var CDR = new(CDR)
	trueResp := Sanitize(event)
	err := json.Unmarshal(trueResp, &CDR)
	if err != nil {
		log.Printf("Cannot unmarshal from event30012: %v\n", err)
		return err
	}

	fmt.Printf("CDR report: %v\n", CDR)
	return nil
}

// (30013) Call Transfer
type CallTransferRequest struct {
	Type            string `db:"type" json:"type"`
	ChannelID       string `db:"channel_id" json:"channel_id"`
	Number          string `db:"number" json:"number"`
	DialPermission  string `db:"dial_permission" json:"dial_permission"`
	AttendedOperate string `db:"attended_operate" json:"attended_operate"`
}

type CallTransferResponse struct {
	Errcode int    `db:"errcode" json:"errcode"`
	Errmsg  string `db:"errmsg" json:"errmsg"`
	CallID  string `db:"call_id" json:"call_id"`
}

// (30014) Call Forward
type CallForwardCallInfo struct {
	CallID  string          `db:"call_id" json:"call_id"`
	Members CallFailMembers `db:"members" json:"members"`
}
type CallForwardReport struct {
	SN  string              `db:"sn" json:"sn"`
	Msg CallForwardCallInfo `db:"msg" json:"msg"`
}
type CallForwardOutboundInfo struct {
	From         string `db:"from" json:"from"`
	To           string `db:"to" json:"to"`
	TrunkName    string `db:"trunk_name" json:"trunk_name"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallForwardInboundInfo struct {
	From         string `db:"from" json:"from"`
	To           string `db:"to" json:"to"`
	TrunkName    string `db:"trunk_name" json:"trunk_name"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallForwardExtensionInfo struct {
	Number       string `db:"number" json:"number"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallForwardMembers struct {
	Extention CallFailExtensionInfo `db:"extention" json:"extention"`
	Inbound   CallFailInboundInfo   `db:"inbound" json:"inbound"`
	Outbound  CallFailOutboundInfo  `db:"outbound" json:"outbound"`
}

// (30015) Call Failed
type CallFailOutboundInfo struct {
	From         string `db:"from" json:"from"`
	To           string `db:"to" json:"to"`
	TrunkName    string `db:"trunk_name" json:"trunk_name"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallFailInboundInfo struct {
	From         string `db:"from" json:"from"`
	To           string `db:"to" json:"to"`
	TrunkName    string `db:"trunk_name" json:"trunk_name"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallFailExtensionInfo struct {
	Number       string `db:"number" json:"number"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
	CallPath     string `db:"call_path" json:"call_path"`
}
type CallFailMembers struct {
	Extention CallFailExtensionInfo `db:"extention" json:"extention"`
	Inbound   CallFailInboundInfo   `db:"inbound" json:"inbound"`
	Outbound  CallFailOutboundInfo  `db:"outbound" json:"outbound"`
}
type CallFailCallInfo struct {
	CallID  string          `db:"call_id" json:"call_id"`
	Reason  string          `db:"reason," json:"reason"`
	Members CallFailMembers `db:"members,omitempty" json:"members,omitempty"`
}

type CallFailedReport struct {
	Type int              `db:"type" json:"type"`
	SN   string           `db:"sn" json:"sn"`
	Msg  CallFailCallInfo `db:"msg" json:"msg"`
}

var CallFailMemberStatuses = map[string]string{
	"ALERT":    "The caller who initiate the call is in the ringback state.",
	"RING":     "The callee is in the ringing state.",
	"ANSWERED": "The call initiated by the caller has been answered.",
	"ANSWER":   "The callee has answered the call, and is in the talking state.",
	"HOLD":     "The call is held.",
	"BYE":      "The call is hung up.",
}

var CallFailReasons = map[string]string{
	"NO Dial Permission":          "The extension/organization has no dial permission.",
	"NO Outbound Restrictin":      "The extension has no outbound call permission.",
	"Circuit/channel congestion":  "Engaged. The channel is in use.",
	"DND":                         "The callee enabled DND.",
	"Line Unreachable":            "The external line is unreachable.",
	"User Busy":                   "The callee declined the call.",
	"410 Gone":                    "Away.",
	"404 NOT found":               "The callee number is not found.",
	"408 Request Time":            "The call is not answered or the callee powers off.",
	"480 Temporatily Unavailable": "The call is not answered.",
	"484 Address incomplete":      "Incorrect format of callee number.",
	"486 Busy here":               "The callee is busy in a call.",
	"603 Declined":                "The call is timed out.",
}

// Handle30015(event) call failed
func (e *CallFailedReport) Handle(event []byte) {
	fmt.Println("Handling event30015")
	var (
		callFailedReport      = new(CallFailedReport)
		callFailCallInfo      = new(CallFailCallInfo)
		callFailMemebersEmbed = new(CallFailMembers)
		callFailExtention     = new(CallFailExtensionInfo)
		callFailInboundInfo   = new(CallFailInboundInfo)
		callFailOutboundInfo  = new(CallFailOutboundInfo)
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
		return
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
		return
	}
	// spew.Dump("Msg: ", messageStruct)
	membersString := messageStruct.Members

	switch membersString {
	case "":
		callFailCallInfo.CallID = messageStruct.CallID
		callFailCallInfo.Reason = messageStruct.Reason
		callFailCallInfo.Members = CallFailMembers{}
		callFailedReport.Type = resp.Type
		callFailedReport.SN = resp.SN
		callFailedReport.Msg = *callFailCallInfo
		spew.Dump(callFailedReport)
		fmt.Println("No memebers found, returning...")
		return
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
			return
		}
		spew.Printf("decoded members: %v of type %T", callFailMemebers)
		// handle Extention
		err = json.Unmarshal([]byte(callFailMemebers.Extention), &callFailExtention)
		if err != nil {
			log.Printf("Cannot unmarshal from Extention: %v\n", err)
			return
		}
		callFailMemebersEmbed.Extention = *callFailExtention
		// handle Inbound
		err = json.Unmarshal([]byte(callFailMemebers.Inbound), &callFailInboundInfo)
		if err != nil {
			log.Printf("Cannot unmarshal from Inbound: %v\n", err)
			return
		}
		callFailMemebersEmbed.Inbound = *callFailInboundInfo
		// handle Outbound
		err = json.Unmarshal([]byte(callFailMemebers.Outbound), &callFailOutboundInfo)
		if err != nil {
			log.Printf("Cannot unmarshal from Outbound: %v\n", err)
			return
		}
		callFailMemebersEmbed.Outbound = *callFailOutboundInfo
	}
	// assign the report values
	callFailCallInfo.Members = *callFailMemebersEmbed
}

// (30016) Inbound Call Invitation
type InboundCallInvitationCallInfoMemebers struct {
	From         string `db:"from" json:"from"`
	To           string `db:"to" json:"to"`
	TrunkName    string `db:"trunk_name" json:"trunk_name"`
	ChannelID    string `db:"channel_id" json:"channel_id"`
	MemberStatus string `db:"member_status" json:"member_status"`
}
type InboundCallInvitationCallInfo struct {
	CallID  string                                `db:"call_id" json:"call_id"`
	Members InboundCallInvitationCallInfoMemebers `db:"members" json:"members"`
}

type InboundCallInvitation struct {
	Type int                           `db:"type" json:"type"`
	SN   string                        `db:"sn" json:"sn"`
	Msg  InboundCallInvitationCallInfo `db:"msg" json:"msg"`
}

// download request
type DownloadRecordingRequest struct {
	ID   int    `db:"id" json:"id"`
	File string `db:"file" json:"file"`
}

// download list response
type DownloadRecordingListResponse struct {
	Errcode             int    `db:"errcode" json:"errcode"`
	Errmsg              string `db:"errmsg" json:"errmsg"`
	File                string `db:"file" json:"file"`
	DownloadResourceURL string `db:"download_resource_url" json:"download_resource_url"`
}

// Query Recording List
type QueryRecordingListRequest struct {
	Page     int    `db:"page" json:"page"`
	PageSize int    `db:"page_size" json:"page_size"`
	SortBy   string `db:"sort_by" json:"sort_by"`
	OrderBy  string `db:"order_by" json:"order_by"`
}

type QueryRecordingListResponse struct { // for later brakedown
	Errcode     int               `db:"errcode" json:"errcode"`
	Errmsg      string            `db:"errmsg" json:"errmsg"`
	TotalNumber int               `db:"total_number" json:"total_number"`
	Data        []RecordingDetail `db:"data" json:"data"`
}
type RecordingDetail struct { // for later brakedown
	ID       int    `db:"id" json:"id"`
	Time     string `db:"time" json:"time"`
	UID      string `db:"uid" json:"uid"`
	CallFrom string `db:"call_from" json:"call_from"`
	CallTo   string `db:"call_to" json:"call_to"`
	Duration int    `db:"duration" json:"duration"`
	Size     int    `db:"size" json:"size"`
	CallType string `db:"call_type" json:"call_type"`
	File     string `db:"file" json:"file"`
}
type RecordingListResponse struct { // for decoding
	Errcode     int    `db:"errcode" json:"errcode"`           // Returned error code.0: Succeed. Non-zero value: Failed.
	Errmsg      string `db:"errmsg" json:"errmsg"`             // Returned message. SUCCESS: Succeed. FAILURE: Failed.
	TotalNumber int    `db:"total_number" json:"total_number"` // The total number of call recording.
	Data        []struct {
		ID       int    `db:"id" json:"id"`               // The unique ID of the call recording.
		Time     string `db:"time" json:"time"`           // The time the call was made or received.
		UID      string `db:"uid" json:"uid"`             // The unique ID of the CDR for which the recording is proceeded.
		CallFrom string `db:"call_from" json:"call_from"` // Caller.
		CallTo   string `db:"call_to" json:"call_to"`     // Callee.
		Duration int    `db:"duration" json:"duration"`   // Call duration.
		Size     int    `db:"size" json:"size"`           // The size of the call recording file. (Unit: Byte)
		CallType string `db:"call_type" json:"call_type"` // Communication type. Inbound Outbound Internal
		File     string `db:"file" json:"file"`           // The name of the call recording file.
	} `db:"data" json:"data"`
}

// Sanitize takes the quotes out of a event response
func Sanitize(event []byte) []byte {
	// "msg":"
	beginBefore := []byte(`"msg":"`)
	beginEnd := []byte(`"msg":`)
	// }"}
	endBefore := []byte(`}"}`)
	endEnd := []byte(`}}`)
	slash := []byte(`\`)
	squat := []byte("")
	newEvent := bytes.Replace(event, beginBefore, beginEnd, 1)
	newEvent = bytes.Replace(newEvent, endBefore, endEnd, 1)
	newEvent = bytes.Replace(newEvent, slash, squat, -1)
	// spew.Dump("sanitized event", newEvent)
	return newEvent
}
