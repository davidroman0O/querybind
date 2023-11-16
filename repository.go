package querybind

import "net/url"

// HTTPContext abstracts the HTTP request and response context.
type HTTPContext interface {
	GetHeader(name string) string
	SetHeader(name, value string)
	QueryParams() url.Values
	OriginalURL() string
}
