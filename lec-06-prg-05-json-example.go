package main

import (
	"encoding/json"
	"fmt"
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
	// Go struct 인스턴스 생성 (Python 딕셔너리 superHeroes_source에 해당)
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

	// Python의 json.dumps(..., indent=4)에 해당하는 기능 수행
	// Go 객체를 들여쓰기 4칸("    ")을 포함한 JSON 바이트 슬라이스로 변환합니다.
	jsonData, err := json.MarshalIndent(superHeroesSource, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	// Python의 print(superHeroes)에 해당하는 기능 수행: JSON 문자열을 콘솔에 출력
	fmt.Println(string(jsonData))
}
