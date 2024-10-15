package http_close

import (
	"net/http"
)

func CloseResp(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Close
}

func CloseReq(resp *http.Request) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}

func Closeresponse(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}
