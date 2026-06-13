GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

vuln:
	govulncheck -show verbose ./...

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/follow cmd/follow/main.go	

debug:
	go run -mod $(GOMOD) cmd/follow/main.go -verbose http://random.localhost


wasmjs:
	GOOS=js GOARCH=wasm \
		go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags wasmjs \
		-o app/follow_wasm/www/wasm/get-posts.wasm \
		cmd/get-posts-wasmjs/main.go
