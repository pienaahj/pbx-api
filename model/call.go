package model

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
