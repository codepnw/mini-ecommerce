package routes

import (
	"fmt"

	"github.com/codepnw/mini-ecommerce/pkg/config"
	"github.com/codepnw/mini-ecommerce/pkg/database"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(cfg *config.EnvConfig) error {
	db, err := database.ConnectPostgres(cfg.DB)
	if err != nil {
		return err
	}
	_ = db

	router := gin.Default()

	port := fmt.Sprintf(":%d", cfg.APP.Port)
	return router.Run(port)
}
