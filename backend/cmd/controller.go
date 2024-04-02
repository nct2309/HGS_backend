package cmd

import (
	"go-jwt/internal/controller"
	"go-jwt/internal/infrastructure/driver"
	"go-jwt/internal/infrastructure/repository"
	"go-jwt/internal/middleware"
	"go-jwt/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (s server) SetupControllers() {

	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.CORS())

	db := driver.ConnectSqlServerDB()

	// init repository
	userRepo := repository.NewUserRepo(db)

	// init usecase
	userUsecase := usecase.NewUserUsecase(userRepo)

	// init controller
	controller.SetupUserRoutes(s.router, userUsecase)
}

func (s server) CloseSqlServerDB() {
	driver.CloseSqlServerDB()
}
