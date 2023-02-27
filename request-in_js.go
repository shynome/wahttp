package wahttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"syscall/js"
)

var ErrRequestAbortThroughJsSignal = errors.New("request abort through js signal")

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

	ctx := req.Context()
	ctx, cancel := context.WithCancelCause(ctx)
	r.Get("signal").Set("onabort", js.FuncOf(func(this js.Value, args []js.Value) any {
		cancel(ErrRequestAbortThroughJsSignal)
		return nil
	}))
	req = req.WithContext(ctx)

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
