package model

// event topics
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
	Extension string
	Status    string
}

// type CallStatusChanged struct {
// 	Type int
// 	SN   string
// 	Msg  ExtCallStatus
// }

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
type CallStatusMemebers struct {
	Extention CallStatusExtensionInfo `db:"extention" json:"extention"`
	Inbound   CallStatusInboundInfo   `db:"inbound" json:"inbound"`
	Outbound  CallStatusOutboundInfo  `db:"outbound" json:"outbound"`
}
type CallStatusCallInfo struct {
	CallID  string           `db:"call_id" json:"call_id"`
	Members CallFailMemebers `db:"members" json:"members"`
}
type CallStatusChanged struct {
	Type int                `db:"type" json:"type"`
	SN   string             `db:"sn" json:"sn"`
	Msg  CallStatusCallInfo `db:"msg" json:"msg"`
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
	CallID  string           `db:"call_id" json:"call_id"`
	Members CallFailMemebers `db:"members" json:"members"`
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
type CallForwardMemebers struct {
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
type CallFailMemebers struct {
	Extention CallFailExtensionInfo `db:"extention" json:"extention"`
	Inbound   CallFailInboundInfo   `db:"inbound" json:"inbound"`
	Outbound  CallFailOutboundInfo  `db:"outbound" json:"outbound"`
}
type CallFailCallInfo struct {
	CallID  string           `db:"call_id" json:"call_id"`
	Reason  string           `db:"reason" json:"reason"`
	Members CallFailMemebers `db:"members" json:"members"`
}
type CallFailedReport struct {
	Type int              `db:"type" json:"type"`
	SN   string           `db:"sn" json:"sn"`
	Msg  CallFailCallInfo `db:"msg" json:"msg"`
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
