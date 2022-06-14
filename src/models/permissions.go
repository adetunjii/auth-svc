package models

type Permissions struct {
	Models
	Name string
}

type Activities struct {
	ID       string
	URL      string
	Method   string
	UrlRegex string
}
