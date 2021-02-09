package href

import "strings"

type Hrefer interface {
	GetHref() string
}

func ParseID(obj Hrefer) string {
	href := obj.GetHref()
	parts := strings.Split(href, "/")
	return parts[len(parts)-1]
}
