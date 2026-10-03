package main

import (
	"log"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/config"
	"github.com/Demiladeolorunsola/devboard-1.0/internal/database"
	"github.com/Demiladeolorunsola/devboard-1.0/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	routes.Setup(router, db, cfg)

	log.Printf("DevBoard listening on port %s", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
