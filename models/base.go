package models

// one return value
type RespCode struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	DevInfo string                 `json:"devinfo"`
	Data    map[string]interface{} `json:"data"`
}

// multi return value
type MrespCode struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	DevInfo string      `json:"devinfo"`
	Data    interface{} `json:"data"`
}

type XRespDetailCode struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type XRespCode struct {
	Error XRespDetailCode `json:"error"`
}

/*
func (rc *RespCode) Error() string {
	return fmt.Sprintf("code: %s, message: %s, data: %v", rc.Code, rc.Message, rc.Data)
}

//
func ErrorResponse(code, message string) *RespCode {
	return &RespCode{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
*/
