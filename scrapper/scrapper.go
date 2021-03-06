package scrapper

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


// Scrape Indeed.kr by a term
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedJob   // jobs는 extractedJob을 요소로 하는 배열
	c := make(chan []extractedJob) // jobs를 보내는 채널
	totalPages := getPages(baseURL) // for문의 범위(length)를 구함
	
	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)  //? getting all the jobs on each page
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
		//? To append the CONTENTS of extractedJobs, simply add '...' => similar to 'Spread Syntax' in JS
	}

	writeJobs(jobs)	//! => the combination of many arrays
	fmt.Println("Done, extracted ", len(jobs))
}


//! 페이지 별 상세내용을 가져오는 함수
//! bridge 역할 (goroutine을 생성해서 job slice 전달받고, main의 channel로 전송)
func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob   // jobs는 extractedJob을 요소로 하는 배열
	c := make(chan extractedJob)
	
	pageURL := url + "&start=" + strconv.Itoa(page * 50) // strconv.Itoa() 는 number => string으로 변환
	fmt.Println(pageURL)

	// GET 요청 보내고, 일단 체크
	res, err := http.Get(pageURL)
	checkErr(err)

	defer res.Body.Close()   //? will be run right after 'getPage func is finished'

	doc, err := goquery.NewDocumentFromReader(res.Body)  // doc은 불러온 html document

	//? Find method => 'jobsearch-SerpJobCard' 이라는 className을 가진 태그를 가져온다
	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)  // channel로 getPage func <-> extractJob func 커뮤니케이션
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c   // 메시지가 channel에 전달되기를 기다렸다가, 메시지를 받으면
		jobs = append(jobs, job)  // jobs라는 배열에 job이라는 변수(extractJob(card))를 하나씩 넣어준다
	}

	mainC <- jobs  // job을 모아둔 배열을 mainChannel에 보냄
}


//! Job 하나를 추출하는 함수
func extractJob(card *goquery.Selection, c chan<-extractedJob) {
	id, _ := card.Attr("data-jk")   //? 'Attr' method는 값, 존재여부를 리턴
	title := CleanString(card.Find(".title > a").Text())  // title class 안의 a 태그를 찾음 => text로 변환		
	location := CleanString(card.Find(".sjcl").Text())
	salary := CleanString(card.Find(".salaryText").Text())
	summary := CleanString(card.Find(".summary").Text())

	//! return할 필요 없음 => channel에 값을 전송하기!
	c <- extractedJob {
		id: id, 
		title: title, 
		location: location, 
		salary: salary, 
		summary: summary,
	}
}


//! 공백을 제거해서 한 줄의 string으로 만들어 주는 함수 (strings 패키지 이용)
func CleanString(str string) string {
	//? TrimSpace로 양쪽에 공백을 없애 줌
	//? -> Fields로 하나의 배열로 만들어 줌
	//? -> Join으로 다시 합쳐줌
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}


//! 페이지가 총 몇개인지 구하는 함수
func getPages(url string) int {
	pages := 0

	// GET 요청 보내고, 일단 체크
	res, err := http.Get(url)
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


//! 가져온 jobs data를 csv파일로 저장해주는 함수
func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)  // check err

	w := csv.NewWriter(file)
	defer w.Flush()  // w 파일 저장 (defer => writeJobs 함수가 끝난 뒤 실행)

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}  // header를 순서대로 정해서 배열에 담음
	wErr := w.Write(headers)  // 배열에 담아놓은 내용을 파일에 입력
	checkErr(wErr)  // check err

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)  // check err
	}
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
