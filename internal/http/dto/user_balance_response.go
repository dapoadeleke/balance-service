package dto

type UserBalanceResponse struct {
	UserID  uint64 `json:"user_id"`
	Balance string `json:"balance"`
}
