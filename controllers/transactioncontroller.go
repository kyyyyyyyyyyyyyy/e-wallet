package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/domain"
	"github.com/kyyyyyyyyyyyyyy/sertifikasibc/models"
)

type BlockChainApi struct {
	Bc *domain.BlockChain
}

func (bca *BlockChainApi) TopUpBalance(ctx *fiber.Ctx) error {
	var request domain.Transaction

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	// Tambahkan validasi logika
	if request.To == "" || request.Amount <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid recipient or amount",
		})
	}

	var user domain.User
	if err := models.DB.Where("public_key = ?", request.To).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	isSuccess := bca.Bc.ToUpBalance(request.To, request.Amount)
	if !isSuccess {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Top-up failed",
		})
	}

	if len(bca.Bc.Pool) == 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Top up succeeded, but no transaction to mine",
		})
	}

	bca.Bc.MineBlock()

	return ctx.JSON(fiber.Map{
		"message": "Top up success, block mined",
	})
}

func (bca *BlockChainApi) Transfer(ctx *fiber.Ctx) error {
	var request domain.Transaction

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	isSuccess := bca.Bc.GiveBalance(request.From, request.To, request.Amount)
	if !isSuccess {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Transfer failed",
		})
	}
	// Setelah transaksi berhasil, mine blok baru
	bca.Bc.MineBlock()
	return ctx.JSON(fiber.Map{"message": "Transfer successful, block mined"})
}

func (bca *BlockChainApi) GetBalance(ctx *fiber.Ctx) error {
	var request domain.Transaction

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	balance := bca.Bc.CalculateBalance(request.From)
	if balance == -1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get balance",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"balance": balance,
	})
}
