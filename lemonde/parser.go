package lemonde

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-liveblog/parser"
	"github.com/anaskhan96/soup"
)

type LeMondeParser struct {
	parser.Parser
}

func init() {
	ctx := context.Background()

	for _, s := range LeMondeSchemes {
		err := parser.RegisterParser(ctx, s, NewLeMondeParser)
		if err != nil {
			panic(err)
		}
	}
}

func NewLeMondeParser(ctx context.Context, uri string) (parser.Parser, error) {
	p := &LeMondeParser{}
	return p, nil
}

func (p *LeMondeParser) GetPosts(ctx context.Context, url string) (string, []string, error) {

	rsp, err := soup.Get(url)

	if err != nil {
		return "", nil, fmt.Errorf("Failed to retrieve %s, %w", url, err)
	}

	posts := make([]string, 0)

	doc := soup.HTMLParse(rsp)

	title_el := doc.Find("title")
	title := title_el.FullText()

	title_parts := strings.Split(title, " | ")
	title = title_parts[0]

	sects := doc.FindAll("section")

	for _, s := range sects {

		s_class := s.Attrs()["class"]

		if s_class != "post post__live-container" {
			continue
		}

		divs := s.FindAll("div")

		for _, d := range divs {

			d_class := d.Attrs()["class"]

			if d_class != "content--live" {
				continue
			}

			posts = append(posts, d.FullText())
		}
	}

	return title, posts, nil
}
