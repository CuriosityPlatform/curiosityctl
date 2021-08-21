package compose

import "fmt"

const (
	metaPrefix = "com.curiosityctl"
)

var (
	waitLabel = fmt.Sprintf("%s.wait", metaPrefix)
	bootLabel = fmt.Sprintf("%s.boot", metaPrefix)
)
