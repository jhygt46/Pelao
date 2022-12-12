package main

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type MyHandler struct{}

func (h *MyHandler) requestHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/init":
			ctx.Response.Header.Set("Content-Type", "application/json")
			var b strings.Builder
			b.Grow(10)
			b.Write([]byte{72, 79, 76, 65, 32, 77, 85, 78, 68, 79})
			fmt.Println("SERVER 1")
			fmt.Fprint(ctx, b.String())
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}
}

func main() {

	pass := &MyHandler{}
	fasthttp.ListenAndServe("80", pass.requestHandler)

}
