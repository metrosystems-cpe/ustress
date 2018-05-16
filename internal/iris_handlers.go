package internal

// CheckArguments checks URL arguments
import (
	"fmt"
	"log"

	"github.com/kataras/iris"
)

var (
	uParam         string
	rParam, wParam int
)

func setCORSAlowAll(ctx iris.Context) iris.Context {
	ctx.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Origin")
	ctx.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Origin")
	ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Expose-Headers", "Date")
	return ctx
}

// URLStress it is a iris function
func URLStress(ctx iris.Context) (iris.Context, error) {
	exampleCall := "?url=http://localhost:9090&requests=20&workers=4"
	if uParam = ctx.URLParam("url"); uParam == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		return ctx, fmt.Errorf("missing url parameter\n eg: %s%s", ctx.Path(), exampleCall)
	}

	rParam, err := ctx.URLParamInt("requests")
	if err != nil || rParam <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		return ctx, fmt.Errorf("missing nr of requests parameter\n eg: %s%s", ctx.Path(), exampleCall)
	}

	wParam, err = ctx.URLParamInt("workers")
	if err != nil || wParam <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		return ctx, fmt.Errorf("missing nr of workers parameter\n eg: %s%s", ctx.Path(), exampleCall)
	}
	// @todo handle error
	messages, _ := NewURLStressReport(uParam, rParam, wParam)
	// os.Stdout.Write(messages)

	// fmt.Println(string(messages))
	ctx = setCORSAlowAll(ctx)
	ctx.ContentType("application/json")
	log.Println(string(messages))
	ctx.WriteString(string(messages))
	// ctx.Writef("%v", fmt.Sprintf("%v", messages))
	// ctx.JSON(fmt.Sprintf("%v", messages))
	return ctx, nil
}
