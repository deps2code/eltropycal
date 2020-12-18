package response

//Response -
type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Data struct {
	Message string `json:"msg"`
}
