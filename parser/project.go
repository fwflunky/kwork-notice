package parser

import (
	"encoding/json"
	"os"
)

type Project struct {
	ID          string
	Link        string
	Title       string
	ReportCount string
	WhatPageWas int
}

func getBaseData() []byte {
	if d, err := os.ReadFile("database.json"); err != nil {
		d = []byte("[]")
		_ = os.WriteFile("database.json", d, os.ModePerm)
		return d
	} else {
		return d
	}
}

func saveBaseData(ss []string) {
	d, _ := json.Marshal(ss)
	_ = os.WriteFile("database.json", d, os.ModePerm)
}

func (p Project) MarkAsSeen() {
	var strct []string
	_ = json.Unmarshal(getBaseData(), &strct)

	saveBaseData(append(strct, p.ID))
}

func (p Project) IsSeen() bool {
	var strct []string
	_ = json.Unmarshal(getBaseData(), &strct)
	for _, s := range strct {
		if s == p.ID {
			return true
		}
	}
	return false
}

func GetOnlyNotSeenProjects(pp []Project) []Project {
	var strct []string
	var out []Project
	mm := map[string]bool{}

	_ = json.Unmarshal(getBaseData(), &strct)
	for _, s := range strct {
		mm[s] = true
	}
	for _, pd := range pp {
		if _, ok := mm[pd.ID]; !ok {
			out = append(out, pd)
		}
	}
	return out
}
