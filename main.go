package main

import (
	"github.com/jess-yon/job-scrapper/scrapper"
	"github.com/labstack/echo/v4"
)

func handleHome(c echo.Context) error {
	// return c.String(http.StatusOK, "Hello, World!") //? 문자열을 전달
	return c.File("home.html")  //?  file 전달
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.Logger.Fatal(e.Start(":1323"))  // server on port 1323

	scrapper.Scrape("term")
}