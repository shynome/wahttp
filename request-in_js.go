package wahttp

import (
	"io"
	"net/http"
	"syscall/js"
)

func JsRequest(r js.Value) (*http.Request, error) {

	var body io.Reader
	jsBody := r.Get("body")
	if !jsBody.IsNull() {
		body = NewJSReader(jsBody.Call("getReader"))
	}

	req, err := http.NewRequest(
		r.Get("method").String(),
		r.Get("url").String(),
		body,
	)
	if err != nil {
		return nil, err
	}

	headersIt := r.Get("headers").Call("entries")
	for {
		e := headersIt.Call("next")
		if e.Get("done").Bool() {
			break
		}
		v := e.Get("value")
		req.Header.Set(v.Index(0).String(), v.Index(1).String())
	}

	return req, nil
}
