package main

import (
	"context"
	"log"

	"github.com/aaronland/go-liveblog/app/follow"
	_ "github.com/aaronland/go-liveblog/guardian"
	_ "github.com/aaronland/go-liveblog/lemonde"
)

func main() {

	ctx := context.Background()
	err := follow.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
