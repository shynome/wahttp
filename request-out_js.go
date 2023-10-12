package wahttp

import (
	"context"
	"net/http"
	"syscall/js"
)

const jsFetchMode = "js.fetch:mode"
const jsFetchCreds = "js.fetch:credentials"
const jsFetchRedirect = "js.fetch:redirect"

func GoRequest(req *http.Request) js.Value {
	rInit := js.Global().Get("Object").New()

	if req.Body != nil && req.ContentLength != 0 {
		jsBody := NewGoReader(req.Body)
		rInit.Set("body", jsBody.Export())
		// https://github.com/nodejs/node/issues/46221#issuecomment-1482439958
		rInit.Set("duplex", "half")
	}

	rInit.Set("method", req.Method)
	if mode := req.Header.Get(jsFetchMode); mode != "" {
		rInit.Set("mode", mode)
		req.Header.Del(jsFetchMode)
	}
	if creds := req.Header.Get(jsFetchCreds); creds != "" {
		rInit.Set("credentials", creds)
		req.Header.Del(jsFetchCreds)
	}
	if redirect := req.Header.Get(jsFetchRedirect); redirect != "" {
		rInit.Set("redirect", redirect)
		req.Header.Del(jsFetchRedirect)
	}

	headersIt := js.Global().Get("Headers").New()
	for k, vv := range req.Header {
		for _, v := range vv {
			headersIt.Call("append", k, v)
		}
	}
	rInit.Set("headers", headersIt)

	ac := js.Global().Get("AbortController").New()
	go func() {
		ctx := req.Context()
		<-ctx.Done()
		reason := context.Cause(ctx).Error()
		ac.Call("abort", reason)
	}()
	rInit.Set("signal", ac.Get("signal"))

	link := req.Host + req.RequestURI
	if req.TLS != nil {
		link = "https://" + link
	} else {
		link = "http://" + link
	}
	r := js.Global().Get("Request").New(link, rInit)
	return r
}
