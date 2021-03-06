// +build appengine

package main

import (
    "net/http"

    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
)

func createMux() *echo.Echo {
    e := echo.New()

    s := standard.New("")
    s.SetHandler(e)
    http.Handle("/", s)

    return e
}