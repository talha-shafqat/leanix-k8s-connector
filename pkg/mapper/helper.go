package mapper

import (
	"strings"
)

var replacer = strings.NewReplacer(
	"/", "_",
	".", "_",
	"+", "_",
	"-", "_",
	"*", "_",
	"\\", "_",
)
