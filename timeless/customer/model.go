package customer

import (
	"github.com/hidenkeys/timeless/room"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	FirstName        *string        `json:"firstName" validate:"required"`
	LastName         *string        `json:"lastName" validate:"required"`
	Phone            *string        `json:"phone" validate:"required"`
	Address          *string        `json:"address"`
	EmergencyContact *string        `json:"emergencyContact"`
	Email            *string        `json:"email" validate:"required,email"`
	PlateNumber      *string        `json:"plateNumber" validate:"required"`
	Bookings         []room.Booking `json:"bookings" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
