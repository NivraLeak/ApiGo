package main

import (
	"encoding/json"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MetaData interface{}

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Synt     []string `json:"synt"`
	Infected string   `json:"infected"`
}

func (u *User) ToJson() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) toPercentageTree() {
	u.Infected = "oh No"
}

type DataSymptom struct {
	data [][]string
}

func (dSymptom *DataSymptom) buildTreeDecision() DecisionNode {
	return buildTreeFor(dSymptom.data)
}

func (dSymptom *DataSymptom) addRow(row []string) string {
	dSymptom.data = append(dSymptom.data, row)
	longRows := len(dSymptom.data)
	longColumns := len(dSymptom.data[0])
	if longRows > 0 {
		dSymptom.buildTreeDecision()
	}

	return dSymptom.data[longRows-1][longColumns-1]

}

func (dSymptom *DataSymptom) buildData() {
	var data = readDataFunc("assets/data/datatest02Covid.csv", 21)

	dSymptom.data = data
}
