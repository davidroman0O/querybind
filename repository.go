package querybind

import "net/url"

// HTTPContext abstracts the HTTP request and response context.
type HTTPContext interface {
	GetHeader(name string) string
	SetHeader(name, value string)
	OriginalURL() string
	Referer() string
	BaseUrl() string
	Path() string
	QueryParams(params *url.Values)
}
