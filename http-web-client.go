package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// performRequest 함수: HTTP 요청을 수행하고 응답을 출력합니다.
func performRequest(method, urlStr, responseLabel string, body io.Reader) {
	fmt.Printf("## %s request for %s\n", method, urlStr)

	// HTTP 클라이언트 생성
	client := &http.Client{}
	
	// 요청 객체 생성
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// 요청 수행
	resp, err := client.Do(req)
	if err != nil {
		// 서버가 실행되고 있지 않다면 연결 오류가 발생할 수 있습니다.
		log.Printf("Error sending %s request to %s: %v", method, urlStr, err)
		return
	}
	defer resp.Body.Close() // 응답 본문을 닫습니다.

	// 응답 본문 읽기
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// 응답 출력
	fmt.Printf("## %s response [start]\n", responseLabel)
	fmt.Println(string(responseBody))
	fmt.Printf("## %s response [end]\n\n", responseLabel)
}

func main() {
	fmt.Println("## HTTP client started.")

	serverURL := "http://localhost:8080"

	// --- 1. GET request for directory retrieval ---
	// 파이썬: requests.get('http://localhost:8080/temp/')
	performRequest("GET", serverURL+"/temp/", "GET", nil)

	// --- 2. GET request for multiplication (query parameters) ---
	// 파이썬: requests.get('http://localhost:8080/?var1=9&var2=9')
	// Go에서는 쿼리 문자열을 URL에 직접 포함하여 요청합니다.
	performRequest("GET", serverURL+"/?var1=9&var2=9", "GET", nil)
	
	// --- 3. POST request for multiplication (form data) ---
	// 파이썬: requests.post('http://localhost:8080', data={'var1':'9','var2':'9'})
	
	// POST 본문 (form-urlencoded) 생성: var1=9&var2=9
	formData := url.Values{}
	formData.Set("var1", "9")
	formData.Set("var2", "9")
	
	// strings.NewReader를 사용하여 데이터를 HTTP 요청 본문으로 전달
	postBody := strings.NewReader(formData.Encode())
	
	fmt.Printf("## POST request for %s with var1 is 9 and var2 is 9\n", serverURL)

	// POST 요청 객체 생성
	req, err := http.NewRequest("POST", serverURL, postBody)
	if err != nil {
		log.Fatalf("Error creating POST request: %v", err)
	}
	
	// Content-Type 헤더 설정 (requests 라이브러리에서 기본으로 설정되던 부분)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	// POST 요청 수행 및 응답 처리
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending POST request: %v", err)
		// 서버가 실행되고 있지 않다면 종료하지 않고 메시지만 출력
	} else {
		defer resp.Body.Close()

		postResponseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading POST response body: %v", err)
		}
		
		fmt.Println("## POST response [start]")
		fmt.Println(string(postResponseBody))
		fmt.Println("## POST response [end]")
	}


	fmt.Println("\n## HTTP client completed.")
}