package www

import (
	"embed"
)

//go:embed index.html css/*.css javascript/*.js wasm/*.wasm
var FS embed.FS
