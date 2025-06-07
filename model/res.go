package model

type Response struct {
	StatusCode int         `json:"status_code"` // mã trạng thái HTTP
	Message    string      `json:"message"`     // thông báo trả về
	Data       interface{} `json:"data"`        // dữ liệu trả về, có thể là bất kỳ kiểu dữ liệu nào
}
