.PHONY: dev
dev:
	npx wrangler dev

.PHONY: build
build:
	templ generate
	go run github.com/syumai/workers/cmd/workers-assets-gen@v0.23.1
	tinygo build -o ./build/app.wasm -target wasm -no-debug ./...

.PHONY: deploy
deploy: build
	npx wrangler deploy
