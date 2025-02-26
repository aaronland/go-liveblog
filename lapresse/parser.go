package lapresse

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/aaronland/go-liveblog/parser"
	"github.com/anaskhan96/soup"
)

type LaPresseParser struct {
	parser.Parser
}

func init() {
	ctx := context.Background()

	for _, s := range LaPresseSchemes {
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

	// <div class="live-center-embed" data-src="https://livecenter.norkon.net/frame/lapresse/60404/default">

	divs := doc.FindAll("div")

	for _, d := range divs {

		d_class := d.Attrs()["class"]
		d_src := d.Attrs()["data-src"]

		if d_class != "live-center-embed" {
			continue
		}

		slog.Debug("Parse Norkon URL", "url", d_src)

		n_posts, err := p.getPostsNorkon(ctx, d_src)

		if err != nil {
			slog.Warn("Failed to derive posts", "url", d_src, "error", err)
			continue
		}

		for _, p := range n_posts {
			posts = append(posts, p)
		}
	}

	return title, posts, nil
}

func (p *LaPresseParser) getPostsNorkon(ctx context.Context, url string) ([]string, error) {

	rsp, err := soup.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve %s, %w", url, err)
	}

	posts := make([]string, 0)

	doc := soup.HTMLParse(rsp)

	divs := doc.FindAll("div")

	allowed_classes := []string{
		"ncpost-content",
		"ncpost-comment-content",
		"ncpost-user-comment",
	}

	for _, d := range divs {

		d_class := d.Attrs()["class"]

		if !slices.Contains(allowed_classes, d_class) {
			slog.Debug("Skip", "class", d_class)
			continue
		}

		slog.Debug("Process", "class", d_class)
		posts = append(posts, d.FullText())
	}

	return posts, nil

}
