package main

import (
	"net/http"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
	"github.com/shynome/err0"
	"github.com/shynome/err0/try"
	"github.com/shynome/wahttp"
)

func main() {
	js.Global().Set("GoFetch", GoFetch())
	// js.Global().Get("console").Call("log", js.Global())
	<-make(chan any)
}

func GoFetch() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		p, resolve, reject := promise.New()

		go func() (err error) {
			defer err0.Then(&err, nil, func() {
				reject(err)
			})
			req := try.To1(wahttp.JsRequest(args[0]))
			resp := try.To1(http.DefaultClient.Do(req))
			resolve(wahttp.GoResponse(resp))
			return
		}()

		return p
	})
}
