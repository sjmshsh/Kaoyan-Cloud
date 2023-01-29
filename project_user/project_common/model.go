package project_common

import "net/http"

type BusinessCode int

type Result struct {
	Code BusinessCode `json:"code"`
	Msg  string       `json:"msg"`
	Data any          `json:"data"`
}

func (r *Result) Success(data any) *Result {
	r.Code = http.StatusOK
	r.Msg = "success"
	r.Data = data
	return r
}

func (r *Result) Fail(code BusinessCode, msg string) *Result {
	r.Code = http.StatusOK
	r.Msg = msg
	return r
}
