package wahttp

import (
	"io"
	"syscall/js"
)

type GoReader struct {
	io.Reader
	chunkSize int
	readed    int
}

func NewGoReader(r io.Reader) *GoReader {
	return &GoReader{
		Reader: r,

		chunkSize: 512,
	}
}

func (r *GoReader) Export() js.Value {
	root := js.Global().Get("Object").New()
	root.Set("type", "bytes")
	// root.Set("autoAllocateChunkSize", r.chunkSize)
	// root.Set("start", js.FuncOf(func(this js.Value, args []js.Value) any {
	// 	// c := controller(args[0])
	// 	return nil
	// }))
	root.Set("pull", js.FuncOf(func(this js.Value, args []js.Value) any {
		c := controller(args[0])
		var dst = make([]byte, r.chunkSize)
		offset := r.readed
		n, err := r.Read(dst)

		if err != nil {
			c.close()
			return nil
		}
		r.readed = offset + n
		c.enqueue(dst[:n])
		return nil
	}))
	// root.Set("cancel", js.FuncOf(func(this js.Value, args []js.Value) any {
	// 	return nil
	// }))
	return root
}

type controller js.Value

func (c controller) enqueue(buf []byte) {
	jsBytes := js.Global().Get("Uint8Array").New(len(buf))
	js.CopyBytesToJS(jsBytes, buf)

	jsBuf := jsBytes.Get("buffer")

	chunk := js.Global().Get("DataView").New(jsBuf)

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
func (c controller) desiredSize() int {
	return js.Value(c).Get("desiredSize").Int()
}

func (c controller) byobRequest() js.Value {
	return js.Value(c).Get("byobRequest")
}
