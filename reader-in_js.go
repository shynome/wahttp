package wahttp

import (
	"io"
	"sync"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
)

type JsReader struct {
	jsReader js.Value //

	r io.Reader
	w *io.PipeWriter

	locker sync.Locker
}

var _ io.Reader = (*JsReader)(nil)
var _ io.Closer = (*JsReader)(nil)

func NewJSReader(v js.Value) *JsReader {
	r, w := io.Pipe()
	return &JsReader{
		jsReader: v,

		r: r, w: w,

		locker: &sync.Mutex{},
	}
}

func (jr *JsReader) readFromJS() {
	jr.locker.Lock()
	defer jr.locker.Unlock()

	pp := jr.jsReader.Call("read")
	v, err := promise.Await(pp)
	if err != nil {
		jr.w.CloseWithError(err)
		return
	}
	if v.Get("done").Bool() {
		jr.w.Close()
		return
	}
	var buf = v.Get("value")
	copyFromJS(jr.w, buf)
}

func copyFromJS(w io.Writer, v js.Value) {
	var readed = 0
	var dst [512]byte
	for {
		n := js.CopyBytesToGo(dst[:], v.Call("subarray", readed))
		readed += n
		if n == 0 {
			break
		}
		w.Write(dst[:n])
	}
}

func (jr *JsReader) Read(p []byte) (n int, err error) {
	go jr.readFromJS()
	return jr.r.Read(p)
}

func (jr *JsReader) Close() error {
	jr.jsReader.Call("cancel")
	return nil
}
