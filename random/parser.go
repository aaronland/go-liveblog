package random

import (
	"context"

	"github.com/aaronland/go-liveblog/parser"
	"github.com/brianvoe/gofakeit/v7"
)

// RandomParser implements the parser.Parser interface by generating
// synthetic data using a fake data library.
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

// NewRandomParser creates a new instance of a RandomParser.
func NewRandomParser(ctx context.Context, uri string) (parser.Parser, error) {
	p := &RandomParser{}
	return p, nil
}

// GetPosts returns a randomly generated title and a single
// paragraph of dummy text. This implementation does not
// perform any actual network requests.
func (p *RandomParser) GetPosts(ctx context.Context, url string) (string, []string, error) {

	return p.GetPostsFromText(ctx, "")
}

func (p *RandomParser) GetPostsFromText(ctx context.Context, txt string) (string, []string, error) {

	title := gofakeit.Sentence(6)
	post := gofakeit.Paragraph()

	return title, []string{post}, nil
}
