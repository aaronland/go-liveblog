package main

import (
	"context"
	"log"

	"github.com/aaronland/go-liveblog/app/follow_wasm"
)

func main() {

	ctx := context.Background()
	err := follow_wasm.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
