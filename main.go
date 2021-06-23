package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=devops&limit=50"

func main() {
	getPages()
}

func getPages() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()   //? will be run right after 'getPages func is finished'

	doc, err := goquery.NewDocumentFromReader(res.Body)  // doc은 불러온 html document
	checkErr(err)

	doc.Find("./pagination")  // pagination 이라는 className을 가진 태그를 가져온다
	
	// fmt.Println(doc)
	return 0  // (temporary return value)
}

// 에러 여부 체크
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// 응답 코드 체크
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}