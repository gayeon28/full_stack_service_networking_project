package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	// 1. 초기 Go struct 인스턴스 생성 (Python의 superHeroes_source에 해당)
	superHeroesSource := Squad{
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
				Powers:         []string{"Radiation resistance", "Turning tiny", "Radiation blast"},
			},
			{
				Name:           "Madame Uppercut",
				Age:            39,
				SecretIdentity: "Jane Wilson",
				Powers:         []string{"Million tonne punch", "Damage resistance", "Superhuman reflexes"},
			},
			{
				Name:           "Eternal Flame",
				Age:            1000000,
				SecretIdentity: "Unknown",
				Powers:         []string{"Immortality", "Heat Immunity", "Inferno", "Teleportation", "Interdimensional travel"},
			},
		},
	}

	// 2. 직렬화 (Marshalling) - Python의 json.dumps() 역할
	// Go struct를 JSON 바이트 슬라이스로 변환합니다. (들여쓰기 4칸 적용)
	jsonData, err := json.MarshalIndent(superHeroesSource, "", "    ")
	if err != nil {
		log.Fatalf("JSON Marshalling 오류: %v", err)
	}
	// jsonData는 Python의 superHeroes_mid에 해당

	// 3. 역직렬화 (Unmarshalling) - Python의 json.loads() 역할
	// JSON 바이트 슬라이스를 새로운 Go struct(superHeroes)로 변환합니다.
	var superHeroes Squad
	err = json.Unmarshal(jsonData, &superHeroes)
	if err != nil {
		log.Fatalf("JSON Unmarshalling 오류: %v", err)
	}
	// superHeroes는 Python의 superHeroes에 해당

	// 4. 데이터 접근 및 출력 (Python의 print(superHeroes['homeTown']) 역할)
	fmt.Println(superHeroes.HomeTown)
}