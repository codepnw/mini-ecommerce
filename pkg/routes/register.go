package routes

import (
	"database/sql"
	"fmt"

	"github.com/codepnw/mini-ecommerce/pkg/config"
	"github.com/codepnw/mini-ecommerce/pkg/database"
	"github.com/codepnw/mini-ecommerce/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type routeConfig struct {
	router *gin.Engine
	db     *sql.DB
	token  *jwt.JWTToken
	tx     database.TxManager
}

func RegisterRoutes(cfg *config.EnvConfig) error {
	db, err := database.ConnectPostgres(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	router := gin.Default()

	token, err := jwt.InitJWT(cfg.JWT)
	if err != nil {
		return err
	}

	tx := database.InitTransaction(db)

	// Register Routes
	routeCfg := &routeConfig{
		router: router,
		db:     db,
		token:  token,
		tx:     tx,
	}
	if err = routeCfg.UserRoutes(); err != nil {
		return err
	}
	routeCfg.ProductRoutes()

	port := fmt.Sprintf(":%d", cfg.APP.Port)
	return router.Run(port)
}
