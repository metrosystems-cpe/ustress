package main

import (
	"github.com/kataras/iris"

	"git.metrosystems.net/reliability-engineering/traffic-monkey/internal"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func newApp() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")

	})
	app.Get("/probe", func(ctx iris.Context) {
		ctx, err := internal.URLStress(ctx)
		if err != nil {
			ctx.Writef(err.Error())
		}
	})
	return app
}

func main() {

	app := newApp()
	app.Run(iris.Addr(":9090"), iris.WithoutServerError(iris.ErrServerClosed))
}
