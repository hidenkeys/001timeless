package room

import (
	"gorm.io/gorm"
	"time"
)

type Booking struct {
	gorm.Model
	CustomerID      *uint           `json:"customerID" validate:"required"`
	Receptionist    uint            `json:"receptionist"`
	Amount          *float64        `json:"amount"`
	IsPaid          bool            `json:"isPaid"`
	PaymentMethod   string          `json:"paymentMethod" validate:"required"`
	IsComplementary bool            `json:"isComplementary" gorm:"default:false"`
	RoomBookings    []*RoomBookings `json:"roomBookings" gorm:"constraint:OnUpdate:CASCADE,onDelete:CASCADE"`
}

type RoomBookings struct {
	gorm.Model
	NumberOfNights uint      `json:"numberOfNights" validate:"required"`
	CheckedIn      bool      `json:"checkedIn" gorm:"default:false"`
	CheckedOut     bool      `json:"checkedOut" gorm:"default:false"`
	StartDate      time.Time `json:"startDate" validate:"required"`
	EndDate        time.Time `json:"endDate" validate:"required"`
	Amount         *float64  `json:"amount"`
	BookingID      uint      `json:"bookingID"`
	RoomID         uint      `json:"roomID"`
}

type Room struct {
	gorm.Model
	Name         *string        `json:"name" validate:"required"`
	Category     *string        `json:"category"`
	Description  *string        `json:"description"`
	Price        float64        `json:"price" validate:"required"`
	Status       *string        `json:"status" gorm:"default:available"`
	RoomBookings []RoomBookings `json:"roomBookings"`
}
