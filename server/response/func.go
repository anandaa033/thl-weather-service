package response

import (
	"context"
	"fmt"
	"thlAttractionService/pkg/config"
	"thlAttractionService/pkg/utility/errorMessage"
)

var allowError = false

type Response struct {
	Code    int         `json:"code,omitempty"`
	Massage string      `json:"massage,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Init(cfg *config.Config) {
	allowError = cfg.App.Config.Allows.Response.Error
}

func New(ctx context.Context, errMsg errorMessage.ErrorMessage, result interface{}, e error) *Response {
	l := fmt.Sprintf("%v", ctx.Value("lang"))
	r := &Response{
		Code:    errMsg.Code,
		Massage: errMsg.Message.Language(l),
		Result:  result,
	}
	if allowError {
		if e != nil {
			r.Error = e.Error()
		}
	}

	return r
}
