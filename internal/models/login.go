package models

type Login struct {
	Username string `json:"username,omitempty"`
	Password []byte `json:"password,omitempty"`
}
