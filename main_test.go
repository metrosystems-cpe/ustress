package main

import (
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestOperationOk(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	e.GET("/probe").WithQueryString("url=http://localhost:9090&requests=20&workers=4").Expect().Status(httptest.StatusOK)
}

// $ go test -v
func TestOperationFail(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	e.GET("/").Expect().Status(httptest.StatusOK)
	e.GET("/probe").Expect().Status(httptest.StatusBadRequest)
	// test missing values
	e.GET("/probe").WithQueryString("url=http://localhost:9090").Expect().Status(httptest.StatusBadRequest)
	e.GET("/probe").WithQueryString("url=http://localhost:9090&requests=20").Expect().Status(httptest.StatusBadRequest)
	// test negative values
	e.GET("/probe").WithQueryString("url=http://localhost:9090&requests=-20&workers=4").Expect().Status(httptest.StatusBadRequest)
	e.GET("/probe").WithQueryString("url=http://localhost:9090&requests=20&workers=-4").Expect().Status(httptest.StatusBadRequest)

}
