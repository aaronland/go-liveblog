package guardian

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-liveblog/parser"
	"github.com/anaskhan96/soup"
)

type GuardianParser struct {
	parser.Parser
}

func init() {
	ctx := context.Background()

	for _, s := range GuardianSchemes {

		err := parser.RegisterParser(ctx, s, NewGuardianParser)

		if err != nil {
			panic(err)
		}
	}
}

func NewGuardianParser(ctx context.Context, uri string) (parser.Parser, error) {
	p := &GuardianParser{}
	return p, nil
}

func (p *GuardianParser) GetPosts(ctx context.Context, url string) (string, []string, error) {

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

	paras := doc.FindAll("p")

	for _, p := range paras {

		p_class := p.Attrs()["class"]

		if !strings.HasPrefix(p_class, "dcr-") {
			continue
		}

		posts = append(posts, p.FullText())
	}

	return title, posts, nil

}
