package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	var jobs []extractedJob   // jobs는 extractedJob을 요소로 하는 배열
	totalPages := getPages() // for문의 범위(length)를 구함
	
	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)  //? getting all the jobs on each page
		jobs = append(jobs, extractedJobs...)
		//? To append the CONTENTS of extractedJobs, simply add '...' => similar to 'Spread Syntax' in JS
	}

	writeJobs(jobs)	//! => the combination of many arrays
	fmt.Println("Done, extracted ", len(jobs))
}


//! 가져온 jobs data를 csv파일로 저장해주는 함수
func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)  // check err

	w := csv.NewWriter(file)
	defer w.Flush()  // w 파일 저장 (defer => writeJobs 함수가 끝난 뒤 실행)

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}  // header를 순서대로 정해서 배열에 담음
	wErr := w.Write(headers)  // 배열에 담아놓은 내용을 파일에 입력
	checkErr(wErr)  // check err

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)  // check err
	}
}


//! 페이지 별 상세내용을 가져오는 함수
func getPage(page int) []extractedJob {
	var jobs []extractedJob   // jobs는 extractedJob을 요소로 하는 배열
	pageURL := baseURL + "&start=" + strconv.Itoa(page * 50) // strconv.Itoa() 는 number => string으로 변환
	fmt.Println(pageURL)

	// GET 요청 보내고, 일단 체크
	res, err := http.Get(pageURL)
	checkErr(err)

	defer res.Body.Close()   //? will be run right after 'getPage func is finished'

	doc, err := goquery.NewDocumentFromReader(res.Body)  // doc은 불러온 html document

	//? Find method => 'jobsearch-SerpJobCard' 이라는 className을 가진 태그를 가져온다
	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)  // jobs라는 배열에 job이라는 변수(extractJob(card))를 하나씩 넣어준다
	})

	return jobs  // job을 모아둔 배열을 return
}


//! Job 하나를 추출하는 함수
func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")   //? 'Attr' method는 값, 존재여부를 리턴
	title := cleanString(card.Find(".title > a").Text())  // title class 안의 a 태그를 찾음 => text로 변환		
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())

	return extractedJob {
		id: id, 
		title: title, 
		location: location, 
		salary: salary, 
		summary: summary,
	}
}


//! 공백을 제거해서 한 줄의 string으로 만들어 주는 함수 (strings 패키지 이용)
func cleanString(str string) string {
	//? TrimSpace로 양쪽에 공백을 없애 줌
	//? -> Fields로 하나의 배열로 만들어 줌
	//? -> Join으로 다시 합쳐줌
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}


//! 페이지가 총 몇개인지 구하는 함수
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


//! 에러 여부를 체크해주는 함수
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}


//! 응답 코드를 체크해주는 함수
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}
