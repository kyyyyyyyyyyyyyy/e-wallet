package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/domain"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/models"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/utils"
	"golang.org/x/crypto/bcrypt"
)

func UserRegister(c *fiber.Ctx) error {
	var payload struct {
		PublicKey string `json:"public_key"`
		Password  string `json:"password"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	// Pastikan NewUser menerima string password, bukan []byte
	user := domain.NewUser(payload.PublicKey, string(hashedPassword), 0)

	// Simpan ke DB
	if err := models.DB.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func UserLogin(c *fiber.Ctx) error {
	var payload struct {
		PublicKey string `json:"public_key"`
		Password  string `json:"password"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	var user domain.User
	if err := models.DB.Where("public_key = ?", payload.PublicKey).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	// Bandingkan password dengan bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentialses",
		})
	}

	// Buat token JWT (opsional, kalau kamu butuh)
	token, err := utils.CreateJWT(user.ID, user.PublicKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}
