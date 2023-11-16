package querybind

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type fiberContext struct {
	*fiber.Ctx
}

func Fiber(ctx *fiber.Ctx) fiberContext {
	return fiberContext{
		Ctx: ctx,
	}
}

func (f fiberContext) GetHeader(name string) string {
	return f.Ctx.Get(name)
}

func (f fiberContext) SetHeader(name, value string) {
	f.Ctx.Set("HX-Push-Url", value)
}

func (f fiberContext) QueryParams(params *url.Values) {
	// Combine with the current HTMX request's query parameters
	f.Ctx.Context().QueryArgs().VisitAll(func(key, value []byte) {
		params.Set(string(key), string(value))
	})
}

func (f fiberContext) OriginalURL() string {
	return string(f.Ctx.Request().Header.Referer())
}

func (f fiberContext) Referer() string {
	return f.Ctx.OriginalURL()
}

func (f fiberContext) BaseUrl() string {
	return f.Ctx.BaseURL()
}

func (f fiberContext) Path() string {
	return f.Ctx.Path()
}
