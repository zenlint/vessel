package web

func Build(ctx *Context) {
	data, err := ctx.Req.Body().Bytes()
	if err != nil {
		ctx.Handle(500, "Fail to receive request body.")
		return
	}

	ctx.Write(data)
}
