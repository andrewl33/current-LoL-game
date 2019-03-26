package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// ChampionCount max number of champions in league of legends
const ChampionCount = 400

// ChampInfo gives info for one champion
type ChampionInfo struct {
	Version string `json:"version"`
	ID      string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Blurb   string `json:"blurb"`
	Info    struct {
		Attack     int `json:"attack"`
		Defense    int `json:"defense"`
		Magic      int `json:"magic"`
		Difficulty int `json:"difficulty"`
	} `json:"info"`
	Image struct {
		Full   string `json:"full"`
		Sprite string `json:"sprite"`
		Group  string `json:"group"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
	} `json:"image"`
	Tags    []string `json:"tags"`
	Partype string   `json:"partype"`
	Stats   struct {
		Hp                   float64 `json:"hp"`
		Hpperlevel           float64 `json:"hpperlevel"`
		Mp                   float64 `json:"mp"`
		Mpperlevel           float64 `json:"mpperlevel"`
		Movespeed            float64 `json:"movespeed"`
		Armor                float64 `json:"armor"`
		Armorperlevel        float64 `json:"armorperlevel"`
		Spellblock           float64 `json:"spellblock"`
		Spellblockperlevel   float64 `json:"spellblockperlevel"`
		Attackrange          float64 `json:"attackrange"`
		Hpregen              float64 `json:"hpregen"`
		Hpregenperlevel      float64 `json:"hpregenperlevel"`
		Mpregen              float64 `json:"mpregen"`
		Mpregenperlevel      float64 `json:"mpregenperlevel"`
		Crit                 float64 `json:"crit"`
		Critperlevel         float64 `json:"critperlevel"`
		Attackdamage         float64 `json:"attackdamage"`
		Attackdamageperlevel float64 `json:"attackdamageperlevel"`
		Attackspeedoffset    float64 `json:"attackspeedoffset"`
		Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
	} `json:"stats"`
}

// AllChampionInfo contains all information from "./champions.json"
type AllChampionInfo struct {
	Type    string                  `json:"type"`
	Format  string                  `json:"format"`
	Version string                  `json:"version"`
	Data    map[string]ChampionInfo `json:"data"`
}

// parses champion json and inserts name into correct array slot
func championJSONtoArray() ([ChampionCount]string, error) {
	var champions [ChampionCount]string
	var championsJSON AllChampionInfo

	jsonFile, err := os.Open("./champions.json")

	if err != nil {
		fmt.Println(err)
		return champions, err
	}
	defer jsonFile.Close()

	byteJSON, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return champions, err
	}

	json.Unmarshal(byteJSON, &championsJSON)

	for name, champInfo := range championsJSON.Data {
		i, err := strconv.Atoi(champInfo.Key)
		if err != nil {
			fmt.Println(err)
			return champions, err
		}
		champions[i] = name
	}
	return champions, nil
}
