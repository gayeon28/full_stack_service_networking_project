package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// simpleCalc: 곱셈 함수 (파이썬의 simple_calc)
func simpleCalc(para1, para2 int) int {
	return para1 * para2
}

// parameterRetrieval: 파라미터 검색 함수 (파이썬의 parameter_retrieval)
func parameterRetrieval(msg string) ([]int, error) {
	result := make([]int, 0, 2)
	fields := strings.Split(msg, "&")
	if len(fields) < 2 {
		return nil, fmt.Errorf("Not enough parameters in message: %s", msg)
	}

	for _, field := range fields {
		parts := strings.Split(field, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid parameter format: %s", field)
		}
		// parts[0] is parameter name (e.g., "var1"), parts[1] is value
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Parameter value is not an integer: %s", parts[1])
		}
		result = append(result, value)
	}
	return result, nil
}

// myHttpHandler: HTTP 요청을 처리하는 핸들러 함수
func myHttpHandler(w http.ResponseWriter, r *http.Request) {
	// 요청 상세 정보 출력 (파이썬의 print_http_request_detail)
	// r.RemoteAddr은 'IP:Port' 형식입니다.
	var clientIP, clientPort string
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// IPv6 또는 포맷이 다를 경우 전체 RemoteAddr을 IP로 사용
		clientIP = r.RemoteAddr
		clientPort = ""
	} else {
		clientIP = host
		clientPort = port
	}

	fmt.Println("::Client address   : ", clientIP)
	fmt.Println("::Client port      : ", clientPort)
	fmt.Println("::Request command  : ", r.Method)
	fmt.Println("::Request line     : ", r.Proto+" "+r.URL.String())
	fmt.Println("::Request path     : ", r.URL.Path)
	fmt.Println("::Request version  : ", r.Proto)

	// 응답 헤더 설정 (파이썬의 send_http_response_header)
	w.Header().Set("Content-Type", "text/html")
	// 상태 코드는 핸들러 내부에서 필요 시 설정하도록 제거함

	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	default:
		// 지원하지 않는 메서드에 대한 응답
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

// handleGet: GET 요청 처리
func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("## handleGet() activated.")

	// 쿼리 파라미터 확인 (http://localhost:8080/?var1=9&var2=9)
	params := r.URL.Query()

	if len(params) > 0 {
		// 계산을 위한 GET 요청
		var1Str := params.Get("var1")
		var2Str := params.Get("var2")

		var1, err1 := strconv.Atoi(var1Str)
		var2, err2 := strconv.Atoi(var2Str)

		if err1 != nil || err2 != nil {
			response := fmt.Sprintf("<html>Error: Invalid parameter(s) for calculation. Expected integers.</html>")
			io.WriteString(w, response)
			fmt.Println("## GET request error: Invalid parameter(s).")
			return
		}

		result := simpleCalc(var1, var2)

		// GET 응답 생성
		response := fmt.Sprintf("<html>GET request for calculation => %d x %d = %d</html>", var1, var2, result)
		io.WriteString(w, response)
		fmt.Printf("## GET request for calculation => %d x %d = %d.\n", var1, var2, result)
	} else {
		// 디렉토리 검색을 위한 GET 요청 (파이썬 코드의 else 블록)
		response := fmt.Sprintf("<html><p>HTTP Request GET for Path: %s</p></html>", r.URL.Path)
		io.WriteString(w, response)
		fmt.Printf("## GET request for directory => %s.\n", r.URL.Path)
	}
}

// handlePost: POST 요청 처리
func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("## handlePost() activated.")

	// 먼저 Body 닫기 예약
	defer r.Body.Close()

	// POST 데이터 읽기
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		log.Printf("Error reading body: %v", err)
		return
	}

	postDataStr := string(body)
	fmt.Printf("## POST request data => %s.\n", postDataStr)

	// URL 쿼리 형식의 POST 데이터 파싱
	values, err := url.ParseQuery(postDataStr)
	if err != nil {
		http.Error(w, "Error parsing POST data", http.StatusBadRequest)
		log.Printf("Error parsing POST data: %v", err)
		return
	}

	var1Str := values.Get("var1")
	var2Str := values.Get("var2")

	var1, err1 := strconv.Atoi(var1Str)
	var2, err2 := strconv.Atoi(var2Str)

	if err1 != nil || err2 != nil {
		response := fmt.Sprintf("<html>Error: Invalid parameter(s) for calculation. Expected integers.</html>")
		io.WriteString(w, response)
		fmt.Println("## POST request error: Invalid parameter(s).")
		return
	}

	result := simpleCalc(var1, var2)

	// POST 응답 생성
	postResponse := fmt.Sprintf("POST request for calculation => %d x %d = %d", var1, var2, result)
	io.WriteString(w, postResponse)
	fmt.Printf("## POST request for calculation => %d x %d = %d.\n", var1, var2, result)
}

func main() {
	serverName := "localhost"
	serverPort := "8080"
	addr := ":" + serverPort

	// http.HandleFunc를 사용하여 모든 경로 "/"에 대해 myHttpHandler 함수를 등록합니다.
	http.HandleFunc("/", myHttpHandler)

	fmt.Printf("## HTTP server started at http://%s:%s.\n", serverName, serverPort)

	// http.ListenAndServe를 사용하여 서버를 시작합니다.
	// 이 함수는 오류가 발생하거나 프로그램이 종료될 때까지 블록됩니다.
	if err := http.ListenAndServe(addr, nil); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("HTTP server stopped.")
		} else {
			log.Fatalf("Error starting server: %v", err)
		}
	}
}
