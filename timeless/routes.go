package main

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hidenkeys/timeless/customer"
	"github.com/hidenkeys/timeless/jwtware"
	"github.com/hidenkeys/timeless/room"
	"github.com/hidenkeys/timeless/user"
)

func requireAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret"), JWTAlg: jwt.SigningMethodHS256.Alg()},
	})
}

func adminOnly(c fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	isAdmin := claims["is_admin"].(bool)

	if !isAdmin {
		return c.Status(http.StatusForbidden).SendString("only admins are permitted")
	}

	return c.Next()
}

// isPaid, customer, employee
func bookingRoutes(r fiber.Router) {
	//r.Use(requireAuth())
	r.Get("/:id", room.GetBookingById)
	r.Patch("/pay", room.ChangePaymentStatus)
	r.Get("/search/getSummary", room.GetBookingSummary)
	r.Post("", room.BookRoom)
	r.Patch("/booking/:bookingId/roomBooking/:roomBookingId", room.UpdateBooking)
	r.Patch("/checkin/:id", room.CheckIn)
	r.Patch("/checkout/:id", room.CheckOut)
	r.Get("/booking/:bookingId/roomBooking/:roomBookingId", room.ViewSingleRoomBooking)
	// extend-stay// get booking by customers
	// export summary

	//r.Use(adminOnly)
	r.Get("", room.GetAllBookings)
	r.Delete("/:id", room.DeleteBooking)
}

func userRoutes(r fiber.Router) {
	// login in
	r.Post("/auth/login", user.Login)
	r.Post("/auth/logout", user.Logout)
	r.Post("", user.CreateEmployee)
	// crud employees --> only admin
	//r.Use(requireAuth())
	r.Get("", user.SearchEmployee) // optional_parameter [role, name, job_role]
	r.Patch("/:id/changePassword", user.ChangePassword)

	//r.Use(adminOnly)
	r.Get("/summary", user.GeneralSummary)
	r.Patch("/:id", user.UpdateEmployee)
	r.Get("/get-all", user.GetAllUsers)
	r.Get("/:id", user.GetById)
	r.Delete("/:id", user.DeleteEmployee)
}

func roomRoutes(r fiber.Router) {
	//r.Use(requireAuth())
	r.Get("", room.SearchWithFilter)
	r.Get("/:id", room.GetById)
	r.Get("/categories", room.GetAllCategories)
	r.Get("/:id/bookedDates", room.GetBookedDates)

	//r.Use(adminOnly)
	r.Post("", adminOnly, room.Create)
	r.Patch("/:id", adminOnly, room.Update)
}

func customerRoutes(r fiber.Router) {
	//r.Use(requireAuth())
	r.Post("", customer.Create)
	r.Get("", customer.GetAll)
	r.Get("/:id", customer.GetById)
	r.Get("/search/findByName", customer.FindByName)
	r.Patch("/:id", customer.Update)
	r.Get("/:id/bookings", customer.GetBookings)

	//r.Use(adminOnly)
	r.Delete("/:id", customer.Delete)
}
