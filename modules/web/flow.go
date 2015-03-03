package web

import (
	"fmt"

	api "github.com/dockercn/anchor"
)

// POST /flow/create
func CreateFlow(ctx *Context, form api.CreateFlowOptions) {
	fmt.Printf("%#v\n", form)
}
