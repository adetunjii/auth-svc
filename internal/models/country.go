package models

type Country struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Iso       string `json:"iso"`
	Name      string `json:"name"`
	PhoneCode int    `json:"phone_code"`
	NiceName  string `json:"nice_name"`
	Iso3      string `json:"iso3"`
	NumCode   int    `json:"num_code"`
	ImageUrl  string `json:"image_url,omitempty"`
}
