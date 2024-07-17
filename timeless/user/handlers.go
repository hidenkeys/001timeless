package user

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hidenkeys/timeless/storage"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length  = 8
	prefix  = "TO-"
)

// generateRandomString generates a random string of fixed length
func generateRandomString(length int, charset string) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

// generateCustomString generates the final string starting with "TO-" followed by 5 random characters
func generateCustomString() (string, error) {
	randomPart, err := generateRandomString(length-len(prefix), charset)
	if err != nil {
		return "", err
	}
	return prefix + randomPart, nil
}

func Login(c fiber.Ctx) error {
	loginRequest := make(map[string]string)
	if err := c.Bind().JSON(&loginRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	var user User
	if result := storage.DB.Raw("SELECT * FROM users WHERE email == @username OR employeeId == @username LIMIT 1", sql.Named("username", loginRequest["username"])).Find(&user); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON("incorect email or password")
	}

	if user.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON("incorect email or password")
	}

	// valid password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest["password"]))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON("incorect email or password")
	}
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": t,
		"user":  user,
	})
}

func Logout(c fiber.Ctx) error {
	// Clear the auth token cookie by setting an expired cookie
	cookie := fiber.Cookie{
		Name:     "authtoken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire the cookie immediately
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func Signup(c fiber.Ctx) error {
	return nil
}

type NewEmployee struct {
	Email            *string `json:"email"`
	Password         string  `json:"password"`
	EmployeeID       *string `json:"employeeID"`
	FirstName        *string `json:"firstName"`
	LastName         *string `json:"lastName"`
	Phone            *string `json:"phone"`
	EmergencyContact *string `json:"emergencyContact"`
	IsAdmin          bool    `json:"isAdmin"`
	Role             string  `json:"role"`
	Salary           float64 `json:"salary"`
}

func CreateEmployee(c fiber.Ctx) error {
	newUser := new(NewEmployee)

	if err := c.Bind().JSON(newUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	result, _ := generateCustomString()
	newUser.EmployeeID = &result

	newDBUser := &User{
		Email:            newUser.Email,
		Password:         newUser.Password,
		EmployeeID:       newUser.EmployeeID,
		FirstName:        newUser.FirstName,
		LastName:         newUser.LastName,
		Phone:            newUser.Phone,
		EmergencyContact: newUser.EmergencyContact,
		IsAdmin:          newUser.IsAdmin,
		Role:             newUser.Role,
		Salary:           newUser.Salary,
	}

	if result := storage.DB.Create(newDBUser); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusCreated).JSON(newUser)
}

func UpdateEmployee(c fiber.Ctx) error {
	newUserInfo := new(User)
	userId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fmt.Errorf("invalid user id"))
	}

	if err = c.Bind().JSON(&newUserInfo); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	user := new(User)
	user.ID = uint(userId)

	if result := storage.DB.Model(user).Updates(newUserInfo); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusOK).JSON(user)
}

func DeleteEmployee(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("invalid customer id")
	}

	if result := storage.DB.Exec("DELETE FROM users WHERE id == ?", id); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.SendStatus(http.StatusNoContent)
}

func SearchEmployee(c fiber.Ctx) error {
	//limit := c.Query("page_size", "10")
	//offset := c.Query("page", "0")
	name := c.Query("name")
	wildcard := "%" + name + "%"

	var users []User

	if result := storage.DB.Raw("SELECT * FROM users where users.firstName like @name or users.lastName like @name or role like @name or email like @name or phone like @name", sql.Named("name", wildcard)).Find(&users); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}
	return c.Status(http.StatusOK).JSON(users)
}

func GetAllUsers(c fiber.Ctx) error {
	//limit := c.Query("page_size", "10")
	//offset := c.Query("page", "0")
	var users []User
	if result := storage.DB.Raw("SELECT * FROM users").Find(&users); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.Status(http.StatusOK).JSON(users)
}

func GetById(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("invalid user id")
	}

	var user User

	if result := storage.DB.Raw("SELECT * FROM users WHERE id == ?", id).Scan(&user); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	if user.ID == 0 {
		return c.Status(http.StatusBadRequest).SendString("invalid user id")
	}

	return c.Status(http.StatusOK).JSON(user)
}

func ChangePassword(c fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).SendString("invalid customer id")
	}

	requestBody := make(map[string]string)

	if err := c.Bind().JSON(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	if requestBody["password"] != requestBody["confirmPassword"] {
		return c.Status(http.StatusInternalServerError).SendString("passwords don't match")
	}

	hashedPassword, err := generateHashPassword(requestBody["password"])
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	if result := storage.DB.Exec("UPDATE users SET password = ? WHERE id == ?", hashedPassword, id); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(result.Error)
	}

	return c.SendStatus(http.StatusOK)
}

func generateHashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), nil
}

func GeneralSummary(c fiber.Ctx) error {
	start := c.Query("start")
	end := c.Query("end")
	var results []struct {
		FirstName      string  `json:"first_name"`
		LastName       string  `json:"last_name"`
		PhoneNumber    string  `json:"phone_number"`
		Address        string  `json:"address"`
		EmailAddress   string  `json:"email_address"`
		PaymentMethod  string  `json:"payment_method"`
		Amount         float64 `json:"amount"`
		CheckinDate    string  `json:"checkin_date"`
		CheckoutDate   string  `json:"checkout_date"`
		NumberOfNights int     `json:"number_of_nights"`
		Receptionist   string  `json:"receptionist"`
		RoomNumber     string  `json:"room_number"`
	}

	err := storage.DB.Raw("SELECT\n    customers.first_name as FirstName,\n    customers.last_name as LastName,\n    customers.phone as PhoneNumber,\n    customers.address as Address,\n    customers.email as EmailAddress,\n    b.payment_method as PaymentMethod,\n    b.amount as Amount,\n    rb.start_date as CheckinDate,\n    rb.end_date as CheckoutDate,\n    number_of_nights as NumberOfNights,\n    b.receptionist as Receptionist,\n    name as RoomNumber\nFROM customers\njoin bookings b on customers.id = b.customer_id\njoin main.room_bookings rb on b.id = rb.booking_id\njoin main.rooms r on rb.room_id = r.id\nwhere (start_date BETWEEN ? AND ? ) AND (end_date BETWEEN ? AND ?)\n", start, end, start, end).Scan(&results).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	f := excelize.NewFile()
	sheet := "Summary"
	f.NewSheet(sheet)

	headers := []string{
		"FirstName", "LastName", "PhoneNumber", "Address", "EmailAddress",
		"PaymentMethod", "Amount", "CheckinDate", "CheckoutDate", "NumberOfNights",
		"Receptionist", "RoomNumber",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	for i, result := range results {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), result.FirstName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), result.LastName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), result.PhoneNumber)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", i+2), result.Address)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", i+2), result.EmailAddress)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", i+2), result.PaymentMethod)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", i+2), result.Amount)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", i+2), result.CheckinDate)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", i+2), result.CheckoutDate)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", i+2), result.NumberOfNights)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", i+2), result.Receptionist)
		f.SetCellValue(sheet, fmt.Sprintf("L%d", i+2), result.RoomNumber)
	}

	filePath := "summary.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate Excel file",
		})
	}

	return c.Download(filePath)
}
