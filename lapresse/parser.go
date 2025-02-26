package lapresse

import (
	"context"
	"fmt"
	"strings"

	"log/slog"
	
	"github.com/aaronland/go-liveblog/parser"
	"github.com/anaskhan96/soup"
)

type LaPresseParser struct {
	parser.Parser
}

func init() {
	ctx := context.Background()

	for _, s := range LaPresseSchemes {

		slog.Info("WUT", "s", s)
		err := parser.RegisterParser(ctx, s, NewLaPresseParser)
		if err != nil {
			panic(err)
		}
	}
}

func NewLaPresseParser(ctx context.Context, uri string) (parser.Parser, error) {
	p := &LaPresseParser{}
	return p, nil
}

func (p *LaPresseParser) GetPosts(ctx context.Context, url string) (string, []string, error) {

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

	slog.Info("DEUG", "title", title)
	
	divs := doc.FindAll("div")

	for _, d := range divs {

		d_class := d.Attrs()["class"]

		slog.Info("WUT", "class", d_class)
		
		if d_class != "ncpost-content" {
			continue
		}

		posts = append(posts, d.FullText())
	}

	return title, posts, nil
}
