//go:build wasmjs
package wasm

import(
	"context"
	"net/url"
	"log/slog"
	"syscall/js"
	"sync"
	"fmt"
	"encoding/json"

	"github.com/aaronland/go-liveblog/parser"	
)

type GetPostsResponse struct {
	Posts []string `json:"posts"`
	Errors []string `json:"errors"`
}

func GetPostsFunc() js.Func{

	// Will this work?
	cache := new(sync.Map)
	
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		source := args[0].String()
		text := args[1].String()
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			ctx := context.Background()

			all_posts := make([]string, 0)
			
			slog.Info("PARSE", "uri", source)
			// START OF put me in sync.Once
			
			u, err := url.Parse(source)
			
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to parse URI, %w", err))
				return nil
			}
			
			p_uri := fmt.Sprintf("%s://", u.Host)
			p, err := parser.NewParser(ctx, p_uri)
			
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to create parser for %s, %w", p_uri, err))
				return nil
			}
			
			// END OF put me in sync.Once
			
			// START OF put me in parser.GetPostsWithCache
			
			_, posts, err := p.GetPostsFromText(ctx, text)
			
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to derive posts for %s, %w", p_uri, err))
				return nil
			}
			
			for _, p := range posts {
				
				_, exists := cache.LoadOrStore(p, true)
				
				if exists {
					continue
				}
				
				all_posts = append(all_posts, p)
			}
			
			// END OF put me in parser.GetPostsWithCache				
			
			rsp := GetPostsResponse{
				Posts: all_posts,
			}
			
			enc, err := json.Marshal(rsp)

			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to marshal posts, %w", err))
				return nil
			}

			resolve.Invoke(string(enc))
			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})

}
	
