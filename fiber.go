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
	// return f.Ctx.get
	return ""
}

func (f fiberContext) SetHeader(name, value string) {

}

func (f fiberContext) QueryParams() url.Values {

}

func (f fiberContext) OriginalURL() string {
	return string(f.Ctx.Request().Header.Referer())
}

func (f fiberContext) Referer() string {
	return f.Ctx.OriginalURL()
}
