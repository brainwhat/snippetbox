package ui

import "embed"

// That's a comment directive telling compliler to store
// files from html/ and static/ in the embedded filesystem `Files`
//
//go:embed "html" "static"
var Files embed.FS
