package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/hidenkeys/timeless/customer"
	"github.com/hidenkeys/timeless/room"
	"github.com/hidenkeys/timeless/storage"
	"github.com/hidenkeys/timeless/user"
)

func main() {
	db, err := storage.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&user.User{}, &room.Booking{}, &room.RoomBookings{}, &customer.Customer{}, &room.Room{})
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(fiber.Config{AppName: "TIMELESS"})

	app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:5173", // Allow requests from this origin
        AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept",
        AllowCredentials: true,
    }))

	api := app.Group("/api/v1")

	bookingsApi := api.Group("/bookings")
	usersApi := api.Group("/users")
	roomsApi := api.Group("/rooms")
	customersApi := api.Group("/customers")

	bookingRoutes(bookingsApi)
	userRoutes(usersApi)
	roomRoutes(roomsApi)
	customerRoutes(customersApi)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
