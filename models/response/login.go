package response

type Login struct {
	AuthToken string `json:"auth_token"`
	LoginAt   string `json:"login_at"`
}
