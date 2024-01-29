package models


type Vendor struct {
	VendorId string `gorm:"primaryKey" json:"vendorId"`
	Name string `json:"name"`
	Contact string `json:"contact"`
	Email string `json:"emailAddress"`
	Address string `json:"address"`
}