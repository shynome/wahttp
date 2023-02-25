package wahttp

import (
	"io"
	"sync"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
)

type JSReader struct {
	jsReader js.Value //

	r io.Reader
	w *io.PipeWriter

	locker sync.Locker
}

func NewJSReader(v js.Value) *JSReader {
	r, w := io.Pipe()
	return &JSReader{
		jsReader: v,

		r: r, w: w,

		locker: &sync.Mutex{},
	}
}

func (jr *JSReader) readFromJS() {
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

func (jr *JSReader) Read(p []byte) (n int, err error) {
	go jr.readFromJS()
	return jr.r.Read(p)
}
