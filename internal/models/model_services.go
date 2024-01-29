package models



type Services struct {
	ServiceId string `gorm:"primaryKey" json:"serviceId"`
	Name string `json:"name"`
	Price float32 `gorm:"type:numeric" json:"price"`
}