package http

import (
	"github.com/dapoadeleke/balance-service/internal/http/dto"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

func (h *Handler) PostTransaction(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(http.StatusBadRequest).SendString("userId is required")
	}

	sourceType := c.Get("Source-Type")
	if sourceType == "" {
		return c.Status(http.StatusBadRequest).SendString("Source-Type header is required")
	}

	var request dto.TransactionRequest
	if err := c.BodyParser(&request); err != nil {
		h.Logger.WithError(err).Error("failed to parse transaction request")
		return c.Status(http.StatusBadRequest).SendString("Invalid request")
	}

	request.UserID = userID
	request.SourceType = sourceType

	if err := request.Validate(); err != nil {
		h.Logger.WithError(err).Error("transaction request validation failed")
		return c.Status(http.StatusBadRequest).SendString("Invalid request")
	}

	transaction, err := request.ToTransaction()
	if err != nil {
		h.Logger.WithError(err).Error("failed to convert request to transaction")
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = h.TransactionService.PostTransaction(c.Context(), transaction)
	if err != nil {
		h.Logger.WithError(err).Error("failed to post transaction")
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).SendString("Transaction successful")
}

func (h *Handler) GetBalance(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	if userIDStr == "" {
		return c.Status(http.StatusBadRequest).SendString("userId is required")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.Logger.WithField("userID", userIDStr).WithError(err).Error("invalid userId")
		return c.Status(http.StatusBadRequest).SendString("Invalid userId")
	}

	response, err := h.UserService.GetBalance(c.Context(), userID)
	if err != nil {
		h.Logger.WithError(err).Error("failed to get user balance")
		return c.Status(http.StatusInternalServerError).SendString("An error occurred")
	}

	return c.Status(http.StatusOK).JSON(response)
}
