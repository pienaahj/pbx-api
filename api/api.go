package api

const (
	https_port string = "8088"
	http_port  string = "80"
	//ip address of the PBX server
	pbx_ip       string = "192.168.5.100"
	Api_path     string = "/openapi/v1.0/"
	Content_type string = "application/json"
)

var (
	// for unsecured connections
	BaseURLUnsecure string = "http://" + pbx_ip + ":" + http_port
	// for secured connections
	BaseURLSecure string = "https://" + pbx_ip + ":" + https_port
	connectionURL string
	access_token  string
	refresh_token string
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
