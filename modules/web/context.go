package web

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/binding"

	api "github.com/dockercn/anchor"

	"github.com/dockercn/vessel/models"
	"github.com/dockercn/vessel/modules/log"
	"github.com/dockercn/vessel/modules/setting"
)

// Context represents context of a request.
type Context struct {
	*macaron.Context
}

type Printable interface {
	PrintName(string) string
}

func (ctx *Context) HasError(form Printable) bool {
	errors := ctx.GetVal(reflect.TypeOf(binding.Errors{})).Interface().(binding.Errors)
	if len(errors) == 0 {
		return false
	}

	err := errors[0]
	switch err.Classification {
	case binding.ERR_CONTENT_TYPE:
		ctx.Handle(422, err.Message)
	case binding.ERR_REQUIRED:
		ctx.Handle(422, form.PrintName(err.Fields()[0])+" can not be empty.")
	}
	fmt.Printf("%#v\n", err)
	return true
}

func (ctx *Context) Handle(status int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Error("%d: %v", status, message)
	if status == http.StatusInternalServerError && setting.ProdMode {
		ctx.JSON(status, map[string]string{
			"message": "Internal server error",
		})
		return
	}
	ctx.JSON(status, map[string]string{
		"message": message,
	})
}

func (ctx *Context) AutoJSON(status int, obj interface{}) {
	switch v := obj.(type) {
	case *models.Flow:
		ctx.JSON(200, &api.Flow{
			UUID:    v.UUID,
			Name:    v.Name,
			Created: v.Created,
		})
		return
	}

	ctx.JSON(status, obj)
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
