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
	js.Global().Set("Fetch", Fetch())
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

func Fetch() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := args[0]
		return js.FuncOf(func(this js.Value, args []js.Value) any {
			p, resolve, reject := promise.New()

			go func() (err error) {
				defer err0.Then(&err, nil, func() {
					reject(err)
				})

				req := try.To1(wahttp.JsRequest(args[0]))
				jsReq := wahttp.GoRequest(req)

				jsResp := handler.Invoke(jsReq)
				jsResp = try.To1(promise.Await(jsResp))
				resp := try.To1(wahttp.JsResponse(jsResp))
				jsResp = wahttp.GoResponse(resp)

				resolve(jsResp)
				return
			}()

			return p
		})

	})
}
