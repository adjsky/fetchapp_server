package ege

type question24Request struct {
	Type int    `json:"type" binding:"required"`
	Char string `json:"char"`
}
