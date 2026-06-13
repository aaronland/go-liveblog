package www

import (
	"embed"
)

//go:embed index.html css/*.css javascript/*.js
var FS embed.FS
