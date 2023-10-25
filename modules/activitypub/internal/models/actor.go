package models

type Actor struct {
	Context           []string  `json:"@context"`
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	PreferredUsername string    `json:"preferred_username"`
	Inbox             string    `json:"inbox"`
	PublicKey         PublicKey `json:"public_key"`
}
type PublicKey struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"public_key_pem"`
}
