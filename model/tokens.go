package model

// Tokens stores token properties that
// are accessed in multiple application layers
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenResponse defines the response from a get_token request
type TokenResponse struct {
	Errcode                   int    `json:"errcode"`
	Errmsg                    string `json:"errmsg"`
	Access_token_expire_time  int    `json:"access_token_expire_time"`
	Access_token              string `json:"access_token"`
	Refresh_token_expire_time int    `json:"refresh_token_expire_time"`
	Refresh_token             string `json:"refresh_token"`
}

type RequestParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
