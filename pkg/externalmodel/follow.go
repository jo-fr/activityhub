package externalmodel

type Activity struct {
	Context      any      `json:"@context"`
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Actor        string   `json:"actor"`
	Object       any      `json:"object"`
	Published    string   `json:"published,omitempty"`
	AttributedTo string   `json:"attributedTo,omitempty"`
	Content      string   `json:"content,omitempty"`
	To           []string `json:"to,omitempty"`
	Sensitive    bool     `json:"sensitive,omitempty"`
}
