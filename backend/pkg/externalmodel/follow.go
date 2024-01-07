package externalmodel

type Activity struct {
	Context      any      `json:"@context"`
	ID           string   `json:"id" `
	Type         string   `json:"type" validate:"required"`
	Actor        string   `json:"actor" validate:"url"`
	Object       any      `json:"object" validate:"required"`
	Published    string   `json:"published,omitempty"`
	AttributedTo string   `json:"attributedTo,omitempty"`
	Content      string   `json:"content,omitempty"`
	To           []string `json:"to,omitempty"`
	Sensitive    *bool    `json:"sensitive,omitempty"` // pointer to allow omitempty
}
