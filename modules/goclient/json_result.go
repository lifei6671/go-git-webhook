package goclient

type JsonResult struct {
	ErrorCode int                 	`json:"error_code"`
	Message string                	`json:"message"`
	Command string			`json:"command,omitempty"`
	MsgId string			`json:"msg_id,omitempty"`
	Data interface{}	      	`json:"data,omitempty"`
}
