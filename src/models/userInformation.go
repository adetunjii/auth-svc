package models

type UserInformation struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address     string `json:"address"`
	State       string `json:"state"`
	Country     string `json:"country"`
	ApartmentNo string `json:"apartment_no"`
}

type UpdateUserInformation struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address     string `json:"address"`
	State       string `json:"state"`
	Country     string `json:"country"`
	ApartmentNo string `json:"apartment_no"`
}
