package wahttp

import (
	"errors"
	"io"
	"sync"
	"syscall/js"
)

type GoReader struct {
	io.Reader
	chunkSize int
	readed    int
	locker    sync.Locker
}

// 16kib
const defaultChunkSize = 16 * 1024

func NewGoReader(r io.Reader) *GoReader {
	return &GoReader{
		Reader: r,

		chunkSize: defaultChunkSize,
		locker:    &sync.Mutex{},
	}
}

func (r *GoReader) SetChunkSize(chunkSize int) {
	r.chunkSize = chunkSize
}

func (r *GoReader) Export() js.Value {
	root := js.Global().Get("Object").New()
	root.Set("type", "bytes")
	root.Set("autoAllocateChunkSize", r.chunkSize)
	// root.Set("start", js.FuncOf(func(this js.Value, args []js.Value) any {
	// 	// c := controller(args[0])
	// 	return nil
	// }))
	root.Set("pull", js.FuncOf(func(this js.Value, args []js.Value) any {
		c := controller(args[0])
		go r.JsRead(c)
		return nil
	}))
	// root.Set("cancel", js.FuncOf(func(this js.Value, args []js.Value) any {
	// 	return nil
	// }))
	rstream := js.Global().Get("ReadableStream").New(root)
	return rstream
}

func (r *GoReader) JsRead(c controller) {
	r.locker.Lock()
	defer r.locker.Unlock()

	var dst = make([]byte, r.chunkSize)
	offset := r.readed
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

	r.readed = offset + n
	c.enqueue(dst[:n])
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
