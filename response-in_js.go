package wahttp

import (
	"net/http"
	"strconv"
	"syscall/js"
)

func JsResponse(r js.Value) (*http.Response, error) {
	resp := &http.Response{
		Status:        r.Get("statusText").String(),
		StatusCode:    r.Get("status").Int(),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		ContentLength: -1,
		Uncompressed:  true,
	}

	if jsBody := r.Get("body"); !jsBody.IsNull() {
		resp.Body = NewJSReader(jsBody.Call("getReader"))
	}

	headersIt := r.Get("headers").Call("entries")
	for {
		e := headersIt.Call("next")
		if e.Get("done").Bool() {
			break
		}
		v := e.Get("value")
		resp.Header.Set(v.Index(0).String(), v.Index(1).String())
	}

	if length := resp.Header.Get("content-length"); length != "" {
		if length, err := strconv.ParseInt(length, 10, 64); err == nil {
			resp.ContentLength = length
		}
	}

	return resp, nil
}
