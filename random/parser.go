package random

import (
	"context"

	"github.com/aaronland/go-liveblog/parser"
	"github.com/brianvoe/gofakeit/v7"
)

type RandomParser struct {
	parser.Parser
}

func init() {
	ctx := context.Background()

	for _, s := range RandomSchemes {

		err := parser.RegisterParser(ctx, s, NewRandomParser)

		if err != nil {
			panic(err)
		}
	}
}

func NewRandomParser(ctx context.Context, uri string) (parser.Parser, error) {
	p := &RandomParser{}
	return p, nil
}

func (p *RandomParser) GetPosts(ctx context.Context, url string) (string, []string, error) {

	title := gofakeit.Sentence(6)
	post := gofakeit.Paragraph()

	return title, []string{post}, nil
}
