package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=devops&limit=50"

func main() {
	totalPages := getPages() // for문의 범위를 구함
	
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}


// 페이지 별 상세내용을 가져오는 함수
func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page * 50) // strconv.Itoa() 는 number => string으로 변환
	fmt.Println(pageURL)
}


// 페이지가 총 몇개인지 구하는 함수
func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()   //? will be run right after 'getPages func is finished'

	doc, err := goquery.NewDocumentFromReader(res.Body)  // doc은 불러온 html document
	checkErr(err)

	//! Find method => 'pagination' 이라는 className을 가진 태그를 가져온다
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()  // <a> tag 가 몇개인지
	})  

	return pages
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