caddy:
	caddy file-server -listen :8000 -root public -browse
wasm:
	cd public && GOOS=js GOARCH=wasm go build -o main.wasm
tinywasm:
	cd public && tinygo build -o tinymain.wasm -panic=trap -no-debug
test: wasm
	deno test -A public/main.ts
