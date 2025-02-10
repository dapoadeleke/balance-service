package http

import (
	"github.com/dapoadeleke/balance-service/internal/db"
	"github.com/dapoadeleke/balance-service/internal/repository"
	"github.com/dapoadeleke/balance-service/internal/service"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Db                 *db.Postgres
	Logger             *log.Logger
	UserService        service.User
	TransactionService service.Transaction
}

func NewHandler(
	db *db.Postgres,
	logger *log.Logger,
) *Handler {

	transactionRepository := repository.NewTransactionRepository(db)
	userRepository := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepository)
	transactionService := service.NewTransactionService(db, userRepository, transactionRepository)

	return &Handler{
		Logger:             logger,
		UserService:        userService,
		TransactionService: transactionService,
	}
}
