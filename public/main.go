package main

import (
	"net/http"
	"syscall/js"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	promise "github.com/nlepage/go-js-promise"
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

		go func() {
			defer err2.Catch(func(err error) {
				reject(err)
			})
			req := try.To1(wahttp.JsRequest(args[0]))
			resp := try.To1(http.DefaultClient.Do(req))
			resolve(wahttp.GoResponse(resp))
		}()

		return p
	})
}
