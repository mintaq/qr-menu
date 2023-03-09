package models

// Response struct to describe response object.

type Response struct {
	Error bool        `json:"error"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// Book ErrorResponse to describe errorResponse object.

type ErrorResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
}

func NewResponse(err bool, msg string, data interface{}) Response {
	return Response{
		Error: err,
		Msg:   msg,
		Data:  data,
	}
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Error: true,
		Msg:   msg,
	}
}

func NewSuccessResponse(data interface{}) Response {
	return Response{
		Error: false,
		Msg:   "success",
		Data:  data,
	}
}
