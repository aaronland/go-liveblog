package main

import (
	"context"
	"log"

	"github.com/aaronland/go-liveblog/app/follow"
	_ "github.com/aaronland/go-liveblog/guardian"
	_ "github.com/aaronland/go-liveblog/lapresse"
	_ "github.com/aaronland/go-liveblog/lemonde"
	_ "github.com/aaronland/go-liveblog/random"
)

func main() {

	ctx := context.Background()
	err := follow.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
