declare global {
  const Go: any;
  const GoFetch: typeof fetch;
}

Deno.test("go-fetch", async () => {
  const p = Deno.run({
    cmd: ["go", "env", "GOROOT"],
    stdout: "piped",
  });
  let GoRoot = await p.output().then((buf) => new TextDecoder().decode(buf));
  GoRoot = GoRoot.slice(0, -1 * "\n".length);

  await import(`${GoRoot}/misc/wasm/wasm_exec.js`);

  const go = new Go();

  const wasmBuf = await Deno.readFile("./public/main.wasm");

  const m = await WebAssembly.instantiate(wasmBuf, go.importObject);

  Promise.resolve().then(() => {
    go.run(m.instance);
  });

  await new Promise((rl) => setTimeout(rl, 0));

  const req = new Request("https://shyno.me");
  const r = await GoFetch(req);

  await Promise.all([r.text(), p.close()]);

  if (r.status != 200) {
    throw new Error(`status is ${r.status}, expect 200`);
  }
});
