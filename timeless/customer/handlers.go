package customer

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/hidenkeys/timeless/room"
	"github.com/hidenkeys/timeless/storage"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Create(c fiber.Ctx) error {
	newCustomer := new(Customer)

	if err := c.Bind().JSON(newCustomer); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	if result := storage.DB.Create(newCustomer); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusCreated).JSON(newCustomer)
}

func Update(c fiber.Ctx) error {
	newCustomerInfo := new(Customer)
	customerId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("invalid customer id"))
	}

	if err = c.Bind().JSON(&newCustomerInfo); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	customer := new(Customer)
	customer.ID = uint(customerId)

	if result := storage.DB.Model(customer).Updates(newCustomerInfo); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusOK).JSON(customer)
}

// FindByName find customer by search {param: [name]}
func FindByName(c fiber.Ctx) error {
	name := c.Query("name", "")

	var customers []Customer

	nameWildCard := "%" + name + "%"

	if result := storage.DB.Raw("SELECT * FROM customers WHERE firstName LIKE @name OR lastName LIKE @name or email like @name or phone like @name or plateNumber like @name", sql.Named("name", nameWildCard)).Find(&customers); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusInternalServerError).JSON(result.Error)
		}
	}

	return c.Status(http.StatusOK).JSON(customers)
}

func Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("invalid customer id")
	}

	if result := storage.DB.Exec("DELETE FROM customers WHERE id == ?", id); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.SendStatus(http.StatusNoContent)
}

func GetBookings(c fiber.Ctx) error {
	//limit := c.Query("page_size", "10")
	//offset := c.Query("page", "0")
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("invalid customer id")
	}

	var bookings []room.Booking
	if result := storage.DB.Preload("RoomBookings").Raw("SELECT * FROM bookings WHERE customer_id == ? ", id).Find(&bookings); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusOK).JSON(bookings)
}

func GetById(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("invalid customer id")
	}

	var customer Customer

	if result := storage.DB.Raw("SELECT * FROM customers WHERE id == ?", id).Scan(&customer); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	if customer.ID == 0 {
		return c.Status(http.StatusBadRequest).SendString("invalid customer id")
	}

	return c.Status(http.StatusOK).JSON(customer)
}

func GetAll(c fiber.Ctx) error {
	//limit := c.Query("page_size", "10")
	//offset := c.Query("page", "0")
	var customers []Customer

	if result := storage.DB.Raw("SELECT * FROM customers").Find(&customers); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusOK).JSON(customers)
}
