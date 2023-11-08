package model

type Activity struct {
	Context string `json:"@context"`
	ID      string `json:"id"`
	Type    string `json:"type"`
	Actor   string `json:"actor"`
	Object  string `json:"object"`
}
