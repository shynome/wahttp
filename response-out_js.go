package wahttp

import (
	"net/http"
	"syscall/js"
)

func GoResponse(resp *http.Response) js.Value {
	jsBody := NewGoReader(resp.Body)
	rInit := js.Global().Get("Object").New()
	rInit.Set("status", resp.StatusCode)
	rInit.Set("statusText", resp.Status)
	headersIt := js.Global().Get("Headers").New()
	for k, vv := range resp.Header {
		for _, v := range vv {
			headersIt.Call("append", k, v)
		}
	}
	rInit.Set("headers", headersIt)
	body := jsBody.Export()
	if resp.StatusCode == http.StatusNoContent {
		body = js.Null()
	}
	r := js.Global().Get("Response").New(body, rInit)
	return r
}
