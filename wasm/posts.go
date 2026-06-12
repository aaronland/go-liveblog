//go:build wasmjs
package wasm

import(
	"context"
	"net/url"
	"syscall/js"
	"sync"
	"fmt"
	"encoding/json"

	"github.com/aaronland/go-liveblog/parser"	
)

func GetPostsFunc() js.Func{

	// Will this work?
	cache := new(sync.Map)
	
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		count := len(args)
		uris := make([]string, count)

		for i := 0; i < count; i++ {
			uris[i] = args[i].String()
		}

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			ctx := context.Background()
			all_posts := make([]string, 0)
			
			for _, uri := range uris {

				// START OF put me in sync.Once
				
				u, err := url.Parse(uri)
				
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
				
				_, posts, err := p.GetPosts(ctx, uri)

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
			}

			enc, err := json.Marshal(all_posts)

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
	
