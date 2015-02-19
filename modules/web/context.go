package web

import (
	"github.com/Unknwon/macaron"
)

// Context represents context of a request.
type Context struct {
	*macaron.Context
}

func (ctx *Context) Handle(status int, message string) {
	ctx.JSON(status, map[string]string{
		"message": message,
	})
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(c *macaron.Context) {
		ctx := &Context{
			Context: c,
		}
		c.Map(ctx)
	}
}
