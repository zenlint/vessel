package web

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/binding"

	api "github.com/containerops/anchor"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/modules/log"
	"github.com/containerops/vessel/modules/setting"
	"github.com/containerops/vessel/modules/utils"
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

func (ctx *Context) Handle(status int, format interface{}, args ...interface{}) {
	var message string
	switch v := format.(type) {
	case error:
		message = v.Error()
	case string:
		message = fmt.Sprintf(v, args...)
	default:
		log.Error("Unknown type: %s", v)
		return
	}
	log.Error("%d: %v", status, message)
	if status == http.StatusInternalServerError && setting.ProdMode {
		ctx.JSON(status, map[string]string{
			"message": "Internal Server Error",
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
			UUID:      v.UUID,
			Name:      v.Name,
			Pipelines: utils.MapToStrings(v.Pipelines),
			Created:   v.Created,
		})
	case *models.Pipeline:
		ctx.JSON(200, &api.Pipeline{
			UUID:     v.UUID,
			Name:     v.Name,
			Requires: utils.MapToStrings(v.Requires),
			Created:  v.Created,
		})
	case *models.Stage:
		ctx.JSON(200, &api.Stage{
			UUID:    v.UUID,
			Name:    v.Name,
			Job:     v.Job,
			Created: v.Created,
		})
	case *models.Job:
		ctx.JSON(200, &api.Job{
			UUID:    v.UUID,
			Name:    v.Name,
			Content: v.Content,
			Created: v.Created,
		})
	default:
		ctx.JSON(status, obj)
	}
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
