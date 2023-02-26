### Intro

export `go-fetch` to js side

### Usage

setup at go wasm side [main.go](public/main.go)

```go
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
```

call at js side [main.ts](public/main.ts)

```ts
import "https://go.dev/misc/wasm/wasm_exec.js";

declare global {
  const Go: any;
  const GoFetch: typeof fetch;
}

Deno.test("go-fetch", async () => {
  const go = new Go();

  const wasmBuf = await Deno.readFile("./public/main.wasm");

  const m = await WebAssembly.instantiate(wasmBuf, go.importObject);

  Promise.resolve().then(() => {
    go.run(m.instance);
  });

  await new Promise((rl) => setTimeout(rl, 0));

  const req = new Request("https://shyno.me");
  const r = await GoFetch(req);

  // clean
  await Promise.all([r.text()]);

  if (r.status != 200) {
    throw new Error(`status is ${r.status}, expect 200`);
  }
});
```
