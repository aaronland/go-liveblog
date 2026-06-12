//go:build wasmjs
package main

import (
	"log/slog"
	"syscall/js"
	
	_ "github.com/aaronland/go-liveblog/guardian"
	_ "github.com/aaronland/go-liveblog/lapresse"
	_ "github.com/aaronland/go-liveblog/lemonde"
	_ "github.com/aaronland/go-liveblog/random"

	"github.com/aaronland/go-liveblog/wasm"
)

func main() {

	posts_func := wasm.GetPostsFunc()
	defer posts_func.Release()

	js.Global().Set("get_posts", posts_func)

	c := make(chan struct{}, 0)

	slog.Info("WASM liveblog functions initialized")
	<-c
}
