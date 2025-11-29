package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type SuperHero struct {
	HomeTown string `json:"homeTown"`
	Active   bool   `json:"active"`
	Members  []Member `json:"members"`
}

type Member struct {
	Powers []string `json:"powers"`
	// 여기에 다른 필드를 추가 가능
}

func main()  {
	fileName := "lec-06-prg-03-json-example.json"

	// 1. 파일 읽기 (Python의 open() + read() 역할)
	// os.ReadFile은 파일을 열고(open), 모든 내용을 읽고(read), 파일을 닫는(close) 작업을 한 번에 수행
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("파일 읽기 오류: %v", err)
		return
	}

	// 2. JSON 디코딩 (Python의 json.load(file) 역할)
	var superHeroes SuperHero
	// Unmarshal은 바이트 슬라이스(JSON 데이터)를 Go 값(&superHeroes)으로 디코딩
	err = json.Unmarshal(data, &superHeroes)
	if err != nil {
		log.Fatalf("JSON 디코딩 오류: %v", err)
		return
	}

	// 3. 데이터 접근 및 출력 (Python 접근 방식과 유사)
	// superHeroes['homeTown'] -> superHeroes.HomeTown
	fmt.Println(superHeroes.HomeTown)

	// superHeroes['active'] -> superHeroes.Active
	fmt.Println(superHeroes.Active)

	// superHeroes['members'][1]['powers'][2] -> superHeroes.Members[1].Powers[2]
	// 배열/슬라이스 접근: [1] (두 번째 요소), [2] (세 번째 요소)
	if len(superHeroes.Members) > 1 && len(superHeroes.Members[1].Powers) > 2 {
		fmt.Println(superHeroes.Members[1].Powers[2])
	} else {
		// 안전을 위해 배열 범위 확인
		fmt.Println("배열 인덱스 범위를 벗어났습니다.")
	}
}