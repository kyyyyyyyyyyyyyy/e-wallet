package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/controllers"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/domain"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/models"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/utils"
)

func main() {
	models.ConnectDB()
	app := fiber.New()

	app.Post("/login", controllers.UserLogin)
	app.Post("/register", controllers.UserRegister)

	// 🔧 Inisialisasi blockchain instance
	blockchain := domain.NewBlockChain()

	// 🔧 Buat instance controller dengan blockchain-nya
	blockchainApi := controllers.BlockChainApi{Bc: blockchain}

	// 🔐 Protected routes
	api := app.Group("/api", utils.JWTProtected())
	api.Post("/topup", blockchainApi.TopUpBalance)
	api.Post("/transfer", blockchainApi.Transfer)
	api.Get("/balance", blockchainApi.GetBalance)

	log.Fatal(app.Listen(":8000"))
}
