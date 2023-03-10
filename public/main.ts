import "https://go.dev/misc/wasm/wasm_exec.js";

declare global {
  const Go: any;
  const GoFetch: typeof fetch;
}

const go = new Go();

const wasmBuf = await Deno.readFile("./public/main.wasm");

const m = await WebAssembly.instantiate(wasmBuf, go.importObject);

Promise.resolve().then(() => {
  go.run(m.instance);
});

Deno.test("go-fetch", async () => {
  const req = new Request("https://shyno.me");
  const r = await GoFetch(req);

  // clean
  await Promise.all([r.text()]);

  if (r.status != 200) {
    throw new Error(`status is ${r.status}, expect 200`);
  }
});

Deno.test("signal abort", async () => {
  // todo
});
