// apps/api-go/cmd/api/main.go (fragmento relevante)
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"turnospro/api-go/internal/auth"
	"turnospro/api-go/internal/db"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../../../.env") // si no está, toma variables del sistema

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ ERROR: no se encontró DB_URL")
	}

	database, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("❌ Error al abrir BD:", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatal("❌ Error al conectar con PostgreSQL:", err)
	}
	fmt.Println("✅ Conectado correctamente a PostgreSQL")

	// (opcional) ejemplo de uso de sqlc para que veas que compila
	_ = db.New(database)
	_ = context.Background()

	// ==== Servidor HTTP ====
	r := gin.Default()

	// Rutas de autenticación
	auth.RegisterRoutes(r, database)

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Servidor activo 🚀"})
	})

	fmt.Println("🌐 Servidor ejecutándose en http://localhost:8080")
	r.Run(":8080")
}
