package models

type Record struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Company     string `json:"company" gorm:"index"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Postal      string `json:"postal"`
	PhoneNumber string `json:"phone_number" gorm:"index"`
	EmailID     string `json:"email_id"`
	WebLink     string `json:"web_link" gorm:"null"`
}
