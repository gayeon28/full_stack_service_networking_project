package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// ResponseBody struct는 서버로부터의 JSON 응답 구조를 나타냅니다.
// 서버의 응답 형식이 { "id": "...", "value": "..." }로 가정합니다.
// Python 코드에서는 { "0001": "..." } 형식이므로, 이를 처리하기 위해 Map을 사용합니다.
type ResponseBody map[string]string

// performRequest 함수: HTTP 요청을 수행하고 응답을 출력하는 범용 함수
func performRequest(step int, method, urlStr string, data url.Values) {
	fmt.Printf("\n#%d %s request to %s\n", step, method, urlStr)

	// 요청 본문 준비
	var body io.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	// HTTP 클라이언트 생성
	client := &http.Client{}

	// 요청 객체 생성
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// POST/PUT 요청 시 Content-Type 설정 (폼 데이터 기준)
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// 요청 수행
	resp, err := client.Do(req)
	if err != nil {
		// 서버 미실행 시 발생 가능
		fmt.Printf("#%d Error sending %s request: %v\n", step, method, err)
		return
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// JSON 디코딩
	var jsonResponse ResponseBody
	if err := json.Unmarshal(responseBodyBytes, &jsonResponse); err != nil {
		// JSON 응답이 아닐 경우
		fmt.Printf("#%d Code: %d >> Raw Body: %s\n", step, resp.StatusCode, string(responseBodyBytes))
		return
	}

	// 결과 출력
	// Python 코드의 r.json()['0001'] 또는 r.json()['0002']에 해당하는 키를 찾기
	key := ""
	if v, ok := data["0001"]; ok && len(v) > 0 {
		key = "0001"
	} else if v, ok := data["0002"]; ok && len(v) > 0 {
		key = "0002"
	} else {
		// URL에서 ID 추출 (GET/DELETE 요청용)
		parts := strings.Split(urlStr, "/")
		if len(parts) > 0 {
			key = parts[len(parts)-1]
		}
	}

	jsonResult := jsonResponse[key]
	
	// Python 출력 형식에 맞춤
	// print("#1 Code:", r.status_code, ">>", "JSON:", r.json(), ">>", "JSON Result:", r.json()['0001'])
	fmt.Printf("#%d Code: %d >> JSON: %v >> JSON Result: %s\n", 
		step, 
		resp.StatusCode, 
		jsonResponse,
		jsonResult,
	)
}

func main() {
	fmt.Println("## Go REST client started.")

	baseURL := "http://127.0.0.1:5000/membership_api/"

	// --- #1 Reads a non registered member : error-case ---
	// r = requests.get('http://127.0.0.1:5000/membership_api/0001')
	performRequest(1, "GET", baseURL+"0001", nil)

	// --- #2 Creates a new registered member : non-error case ---
	// r = requests.post('http://127.0.0.1:5000/membership_api/0001', data={'0001':'apple'})
	formData2 := url.Values{"0001": {"apple"}}
	performRequest(2, "POST", baseURL+"0001", formData2)

	// --- #3 Reads a registered member : non-error case ---
	// r = requests.get('http://127.0.0.1:5000/membership_api/0001')
	performRequest(3, "GET", baseURL+"0001", nil)

	// --- #4 Creates an already registered member : error case ---
	// r = requests.post('http://127.0.0.1:5000/membership_api/0001', data={'0001':'xpple'})
	formData4 := url.Values{"0001": {"xpple"}}
	performRequest(4, "POST", baseURL+"0001", formData4)

	// --- #5 Updates a non registered member : error case ---
	// r = requests.put('http://127.0.0.1:5000/membership_api/0002', data={'0002':'xrange'})
	formData5 := url.Values{"0002": {"xrange"}}
	performRequest(5, "PUT", baseURL+"0002", formData5)

	// --- #6 Updates a registered member : non-error case ---
	// Python 코드에서 두 번째 요청은 POST, 세 번째 요청은 PUT입니다.
	// 1. 등록 (POST)
	formData6_1 := url.Values{"0002": {"xrange"}}
	performRequest(6, "POST", baseURL+"0002", formData6_1)
	
	// 2. 수정 (PUT)
	formData6_2 := url.Values{"0002": {"orange"}}
	performRequest(6, "PUT", baseURL+"0002", formData6_2)


	// --- #7 Delete a registered member : non-error case ---
	// r = requests.delete('http://127.0.0.1:5000/membership_api/0001')
	performRequest(7, "DELETE", baseURL+"0001", nil)

	// --- #8 Delete a non registered member : non-error case ---
	// r = requests.delete('http://127.0.0.1:5000/membership_api/0001')
	performRequest(8, "DELETE", baseURL+"0001", nil)

	fmt.Println("\n## Go REST client completed.")
}