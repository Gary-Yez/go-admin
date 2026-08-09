package templates

import "embed"

// FS contains the code-generation templates distributed with go-admin.
//
//go:embed server/*.tmpl web/*.tmpl
var FS embed.FS
