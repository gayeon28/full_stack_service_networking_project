package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync" // 동시성 제어를 위한 패키지
)

// MembershipHandler: Python의 MembershipHandler 클래스에 해당하는 Go Struct
type MembershipHandler struct {
	// database: 회원 정보를 저장하는 Map (id: value)
	// Mutex: Map 접근 시 동시성 문제를 해결하기 위한 읽기/쓰기 락
	database map[string]string
	mu       sync.RWMutex
}

// 응답 구조체
type Response struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// 새 핸들러 인스턴스를 생성하는 생성자
func NewMembershipHandler() *MembershipHandler {
	return &MembershipHandler{
		database: make(map[string]string),
	}
}

// =================================================================
// Go 언어의 RESTful 핸들러 함수들
// =================================================================

// handleErrorResponse: 오류 응답을 처리하고 클라이언트에게 JSON을 전송
func handleErrorResponse(w http.ResponseWriter, id string, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{ID: id, Value: message}
	json.NewEncoder(w).Encode(response)
}

// handleSuccessResponse: 성공 응답을 처리하고 클라이언트에게 JSON을 전송
func handleSuccessResponse(w http.ResponseWriter, id string, value string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{ID: id, Value: value}
	json.NewEncoder(w).Encode(response)
}

// create (POST): 새 멤버를 추가
func (m *MembershipHandler) create(w http.ResponseWriter, r *http.Request, memberID string) {
	// Python 코드에서는 request.form[member_id]를 사용했습니다.
	// Go에서는 POST Body에서 해당 값을 파싱해야 합니다.
	
	// 폼 데이터 파싱
	r.ParseForm()
	value := r.FormValue("value") // 요청 본문에서 'value' 필드를 추출한다고 가정
	
	if value == "" {
		// Python 코드에서는 폼 키가 member_id와 같았으나, 일반적인 REST 방식에 맞춰 'value'로 가정
		value = r.FormValue(memberID)
	}

	if value == "" {
		handleErrorResponse(w, memberID, "Value field missing in POST data", http.StatusBadRequest)
		return
	}

	// 락 획득 (쓰기)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.database[memberID]; exists {
		handleSuccessResponse(w, memberID, "None", http.StatusOK) // 이미 존재하면 "None" 반환
		return
	}

	m.database[memberID] = value
	handleSuccessResponse(w, memberID, value, http.StatusCreated) // 201 Created
}

// read (GET): 멤버 정보를 조회
func (m *MembershipHandler) read(w http.ResponseWriter, r *http.Request, memberID string) {
	// 락 획득 (읽기)
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.database[memberID]
	if !exists {
		handleSuccessResponse(w, memberID, "None", http.StatusOK)
		return
	}

	handleSuccessResponse(w, memberID, value, http.StatusOK)
}

// update (PUT): 멤버 정보를 수정
func (m *MembershipHandler) update(w http.ResponseWriter, r *http.Request, memberID string) {
	// PUT 요청에서 본문 데이터 읽기 (Go의 PUT은 폼 데이터를 자동으로 파싱하지 않음)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleErrorResponse(w, memberID, "Error reading request body", http.StatusInternalServerError)
		return
	}
	
	// 본문 데이터를 폼 형식(URL 쿼리 형식)으로 파싱
	values := parseFormData(string(body))
	if err != nil {
		handleErrorResponse(w, memberID, "Invalid PUT body format", http.StatusBadRequest)
		return
	}
	
	value := values["value"] // 요청 본문에서 'value' 필드를 추출한다고 가정
	
	// 락 획득 (쓰기)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.database[memberID]; !exists {
		handleSuccessResponse(w, memberID, "None", http.StatusOK) // 존재하지 않으면 "None" 반환
		return
	}

	m.database[memberID] = value
	handleSuccessResponse(w, memberID, value, http.StatusOK)
}

// delete (DELETE): 멤버 정보를 삭제
func (m *MembershipHandler) delete(w http.ResponseWriter, r *http.Request, memberID string) {
	// 락 획득 (쓰기)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.database[memberID]; !exists {
		handleSuccessResponse(w, memberID, "None", http.StatusOK) // 존재하지 않으면 "None" 반환
		return
	}

	delete(m.database, memberID)
	handleSuccessResponse(w, memberID, "Removed", http.StatusOK)
}

// parseFormData: PUT 요청의 본문 데이터를 폼 형식으로 파싱하는 헬퍼 함수
func parseFormData(data string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(data, "&")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			// URL 디코딩을 수행해야 하지만, 여기서는 단순 문자열로 가정
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// =================================================================
// 라우팅 및 메인 함수
// =================================================================

// mainHandler: 모든 HTTP 메서드(POST/GET/PUT/DELETE)를 처리하는 단일 핸들러
func (m *MembershipHandler) mainHandler(w http.ResponseWriter, r *http.Request) {
	// URL 경로에서 member_id 추출 (Go의 net/http는 경로 변수를 자동으로 추출하지 않으므로 수동 파싱)
	// /membership_api/member_id 형식에서 member_id 부분을 추출합니다.
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[len(parts)-2] != "membership_api" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	memberID := parts[len(parts)-1]
	
	// HTTP 메서드에 따라 적절한 CRUD 함수 호출
	switch r.Method {
	case "POST":
		m.create(w, r, memberID)
	case "GET":
		m.read(w, r, memberID)
	case "PUT":
		m.update(w, r, memberID)
	case "DELETE":
		m.delete(w, r, memberID)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func main() {
	// 핸들러 인스턴스 생성
	myManager := NewMembershipHandler()

	// 라우팅 설정: 모든 /membership_api/* 경로 요청을 myManager.mainHandler가 처리하도록 합니다.
	http.HandleFunc("/membership_api/", myManager.mainHandler)

	addr := ":5000" // Flask 기본 포트 5000을 사용
	fmt.Printf("## RESTful API Server started at http://localhost%s\n", addr)
	
	// 서버 시작
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}