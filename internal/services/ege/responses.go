package ege

type questionResponse struct {
	Code   int `json:"code"`
	Result int `json:"result"`
}

type availableResponse struct {
	Code               int    `json:"code"`
	QuestionsAvailable string `json:"questions_available"`
}

type questionTypesResponse struct {
	Code           int    `json:"code"`
	TypesAvailable string `json:"types_available"`
}
