package models

type Response struct {
	Error bool        `json:"error"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func NewResponse(err bool, msg string, data interface{}) Response {
	return Response{
		Error: err,
		Msg:   msg,
		Data:  data,
	}
}
