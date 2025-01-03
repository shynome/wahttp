package wahttp

import (
	"io"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
)

type JsWriter struct {
	jsWriter js.Value //
}

var _ io.Writer = (*JsWriter)(nil)
var _ io.Closer = (*JsWriter)(nil)

func NewJSWriter(v js.Value) *JsWriter {
	return &JsWriter{
		jsWriter: v,
	}
}

var jsUint8Array = js.Global().Get("Uint8Array")

func (jw *JsWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	body := jsUint8Array.New(n)
	js.CopyBytesToJS(body, p)
	pp := jw.jsWriter.Call("write", body)
	_, err = promise.Await(pp)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (jw *JsWriter) Close() error {
	pp := jw.jsWriter.Call("close")
	_, err := promise.Await(pp)
	return err
}
