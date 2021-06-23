package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

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

	// GET 요청 보내고, 일단 체크
	res, err := http.Get(pageURL)
	checkErr(err)

	defer res.Body.Close()   //? will be run right after 'getPage func is finished'

	doc, err := goquery.NewDocumentFromReader(res.Body)  // doc은 불러온 html document

	//! Find method => 'jobsearch-SerpJobCard' 이라는 className을 가진 태그를 가져온다
	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")   //? 'Attr' method는 값, 존재여부를 리턴
		title := card.Find(".title > a").Text()  // title class 안의 a 태그를 찾음 => text로 변환		
		location := card.Find(".sjcl").Text()

		fmt.Println(id, title, location)
	})
}


// 페이지가 총 몇개인지 구하는 함수
func getPages() int {
	pages := 0

	// GET 요청 보내고, 일단 체크
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()   //? will be run right after 'getPages func is finished'

	doc, err := goquery.NewDocumentFromReader(res.Body)  // doc은 불러온 html document
	checkErr(err)

	//! Find method => 'pagination' 이라는 className을 가진 태그를 가져온다
	doc.Find(".pagination").Each(func(i int, card *goquery.Selection) {  // 'card'는 각각의 job link 카드를 가져온 것
		pages = card.Find("a").Length()  // <a> tag 가 몇개인지
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
