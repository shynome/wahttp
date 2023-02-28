package wahttp

import (
	"errors"
	"io"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
)

type GoReader struct {
	io.ReadCloser

	chunkSize int
	chunkBuf  []byte
}

// 16kib
const defaultChunkSize = 16 * 1024

func NewGoReader(r io.ReadCloser) *GoReader {
	return &GoReader{
		ReadCloser: r,

		chunkSize: defaultChunkSize,
	}
}

func (r *GoReader) SetChunkSize(chunkSize int) {
	r.chunkSize = chunkSize
}

func (r *GoReader) Export() js.Value {
	root := js.Global().Get("Object").New()
	root.Set("type", "bytes")
	root.Set("autoAllocateChunkSize", r.chunkSize)
	root.Set("start", js.FuncOf(func(this js.Value, args []js.Value) any {
		r.chunkBuf = make([]byte, r.chunkSize)
		return nil
	}))
	root.Set("pull", js.FuncOf(func(this js.Value, args []js.Value) any {
		c := controller(args[0])
		p, resolve, reject := promise.New()
		go func() {
			if err := r.JsRead(c); err != nil {
				reject(err.Error())
				return
			}
			resolve(1)
		}()
		return p
	}))
	root.Set("cancel", js.FuncOf(func(this js.Value, args []js.Value) any {
		p, resolve, reject := promise.New()
		go func() {
			if err := r.Close(); err != nil {
				reject(err.Error())
				return
			}
			resolve(1)
		}()
		return p
	}))
	rstream := js.Global().Get("ReadableStream").New(root)
	return rstream
}

func (r *GoReader) JsRead(c controller) (err error) {

	var dst = r.chunkBuf
	n, err := r.Read(dst)

	if errors.Is(err, io.EOF) {
		if n == 0 {
			c.close()
			return
		}
		err = nil
	}

	if err != nil {
		c.error(err)
		return
	}

	c.enqueue(dst[:n])

	return
}

type controller js.Value

func (c controller) enqueue(buf []byte) {
	jsBytes := js.Global().Get("Uint8Array").New(cap(buf))
	n := js.CopyBytesToJS(jsBytes, buf)

	jsBuf := jsBytes.Get("buffer")
	chunk := js.Global().Get("DataView").New(jsBuf, 0, n)

	req := c.byobRequest()
	if !req.IsNull() {
		req.Call("respondWithNewView", chunk)
		return
	}

	js.Value(c).Call("enqueue", chunk)
}
func (c controller) close() {
	js.Value(c).Call("close")
}
func (c controller) error(err error) {
	js.Value(c).Call("error", err.Error())
}
func (c controller) desiredSize() int {
	return js.Value(c).Get("desiredSize").Int()
}

func (c controller) byobRequest() js.Value {
	return js.Value(c).Get("byobRequest")
}
