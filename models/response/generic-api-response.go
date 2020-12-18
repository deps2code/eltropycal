package response

//Response -
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type Data struct {
	Message string `json:"msg"`
}
