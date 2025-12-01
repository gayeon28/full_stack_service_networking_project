package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Member struct는 squad 멤버의 세부 정보를 나타냅니다.
type Member struct {
	Name           string   `json:"name"`
	Age            int      `json:"age"`
	SecretIdentity string   `json:"secretIdentity"`
	Powers         []string `json:"powers"`
}

// Squad struct는 전체 슈퍼 히어로 분대 정보를 나타냅니다.
type Squad struct {
	SquadName  string   `json:"squadName"`
	HomeTown   string   `json:"homeTown"`
	Formed     int      `json:"formed"`
	SecretBase string   `json:"secretBase"`
	Active     bool     `json:"active"`
	Members    []Member `json:"members"`
}

func main() {
	// Go struct 인스턴스 생성 (Python 딕셔너리에 해당)
	superHeroes := Squad{
		SquadName:  "Super hero squad",
		HomeTown:   "Metro City",
		Formed:     2016,
		SecretBase: "Super tower",
		Active:     true,
		Members: []Member{
			{
				Name:           "Molecule Man",
				Age:            29,
				SecretIdentity: "Dan Jukes",
				Powers: []string{
					"Radiation resistance",
					"Turning tiny",
					"Radiation blast",
				},
			},
			{
				Name:           "Madame Uppercut",
				Age:            39,
				SecretIdentity: "Jane Wilson",
				Powers: []string{
					"Million tonne punch",
					"Damage resistance",
					"Superhuman reflexes",
				},
			},
			{
				Name:           "Eternal Flame",
				Age:            1000000,
				SecretIdentity: "Unknown",
				Powers: []string{
					"Immortality",
					"Heat Immunity",
					"Inferno",
					"Teleportation",
					"Interdimensional travel",
				},
			},
		},
	}

	// Python 코드의 print 문에 해당하는 필드 접근 및 출력
	fmt.Println(superHeroes.HomeTown) // 'Metro City' 출력
	fmt.Println(superHeroes.Active)   // 'true' 출력
	// Members[1] (Madame Uppercut)의 Powers[2] ("Superhuman reflexes") 접근
	fmt.Println(superHeroes.Members[1].Powers[2])

	// Go 객체를 JSON 바이트로 마샬링 (인코딩)
	// json.MarshalIndent를 사용하여 들여쓰기('\t')를 포함한 예쁜 출력 생성
	jsonData, err := json.MarshalIndent(superHeroes, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	// JSON 데이터를 파일에 쓰기
	fileName := "lec-06-prg-04-json-example.json"
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}

	fmt.Printf("\nJSON data successfully written to %s\n", fileName)
}
