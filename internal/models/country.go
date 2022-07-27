package models

type Country struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Iso  string `json:"iso"`
	Name string `json:"name"`
	//PhoneCode string   `json:"phoneCode"`
	NiceName string `json:"niceName"`
	Iso3     string `json:"iso3"`
	//NumCode   int    `json:"numCode"`
}
