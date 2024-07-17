package user

import (
	"github.com/hidenkeys/timeless/room"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email            *string `json:"email" gorm:"unique" validate:"required,email"`
	Password         string
	EmployeeID       *string `json:"employeeID" gorm:"unique" validate:"required"`
	FirstName        *string `json:"firstName" validate:"required"`
	LastName         *string `json:"lastName" validate:"required"`
	Phone            *string `json:"phone" validate:"required"`
	EmergencyContact *string `json:"emergencyContact"`
	IsAdmin          bool    `json:"isAdmin" gorm:"default:0"`
	Role             string  `json:"role"`
	Salary           float64 `json:"salary"`

	Bookings []room.Booking `json:"bookings" gorm:"foreignKey:Receptionist"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		u.Password, err = generateHashPassword(u.Password)
	} else {
		u.Password, err = generateHashPassword("password")
	}

	return
}
