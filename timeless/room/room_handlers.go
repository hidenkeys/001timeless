package room

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/hidenkeys/timeless/storage"
	"net/http"
	"strconv"
)

func Create(c fiber.Ctx) error {
	newRoom := new(Room)

	if err := c.Bind().JSON(newRoom); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	//status := "available"
	//newRoom.Status = &status

	if result := storage.DB.Create(newRoom); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusCreated).JSON(newRoom)
}

func Update(c fiber.Ctx) error {
	newRoomInfo := make(map[string]any)
	roomID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("invalid room id"))
	}

	if err = c.Bind().JSON(&newRoomInfo); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	room := new(Room)
	room.ID = uint(roomID)

	if result := storage.DB.Model(room).Updates(newRoomInfo); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusOK).JSON(room)
}

// SearchWithFilter params {filter: [name, category, status] , name, category, status }
// getRoomByName ---
// getRoomByCategory ---
// getRoomByStatus ---
func SearchWithFilter(c fiber.Ctx) error {
	filter := c.Query("filter")
	value := c.Query("value")

	var rooms []Room

	switch filter {
	case "":
		if result := storage.DB.Raw("SELECT * FROM rooms").Find(&rooms); result.Error != nil {
			return c.Status(http.StatusInternalServerError).JSON(result.Error)
		} else {
			return c.Status(http.StatusOK).JSON(rooms)
		}
	default:
		if result := storage.DB.Raw(fmt.Sprintf("SELECT * FROM rooms WHERE %s == ?", filter), value).Find(&rooms); result.Error != nil {
			return c.Status(http.StatusInternalServerError).JSON(result.Error)
		} else {
			return c.Status(http.StatusOK).JSON(rooms)
		}
	}
}

func GetById(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("invalid customer id")
	}

	var rooms []Room

	if result := storage.DB.Raw("SELECT * FROM rooms WHERE id == ?", id).Find(&rooms); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusCreated).JSON(rooms)
}

func GetAllCategories(c fiber.Ctx) error {
	var categories []string

	if result := storage.DB.Raw("SELECT DISTINCT category FROM rooms").Find(&categories); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusOK).JSON(categories)
}
