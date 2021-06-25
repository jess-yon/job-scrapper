package main

import (
	"job-scrapper/scrapper"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error {
	// return c.String(http.StatusOK, "Hello, World!") //? 문자열 전달
	return c.File("home.html")  //?  file 전달
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)  // handleScrape func를 통해 사용자에게 file 보낸 후, 서버에서는 파일을 remove

	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	
	// 첨부파일을 return 해주는 기능 (저장된 파일 이름, 전달할 파일 이름)
	return c.Attachment(fileName, fileName)  //! => 'job.scv'라는 파일이 다운로드 됨
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))  // server on port 1323
}
