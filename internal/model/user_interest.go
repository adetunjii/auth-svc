package model

type Interest struct {
	BaseEntity

	Title string `json:"title"`
}

type UserInterest struct {
	BaseEntity

	UserId     string `json:"user_id"`
	InterestId string `json:"interest_id"`
}
