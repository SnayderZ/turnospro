// apps/api-go/cmd/api/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"turnospro/api-go/internal/auth"
	"turnospro/api-go/internal/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// 🧩 1️⃣ Intentar cargar el archivo .env desde varias ubicaciones comunes
	err := godotenv.Load(
		"../../../../.env",
	)
	if err != nil {
		log.Println("⚠️ No se encontró archivo .env, usando variables del sistema...")
	}

	// 🧩 2️⃣ Verificar la variable DB_URL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ ERROR: no se encontró la variable DB_URL (verifica tu .env)")
	}

	// 🧩 3️⃣ Conexión a la base de datos
	database, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("❌ Error al abrir conexión con BD:", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatal("❌ Error al conectar con PostgreSQL:", err)
	}
	fmt.Println("✅ Conectado correctamente a PostgreSQL")

	// 🧩 4️⃣ Inicialización de rutas y servidor
	r := gin.Default()

	// Registrar rutas de autenticación
	auth.RegisterRoutes(r, database)

	// ====== Rutas protegidas 🔒 ======
	protected := r.Group("/api/protected")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", auth.ProfileHandler)
	}

	// Ruta de prueba (health check)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Servidor activo 🚀"})
	})

	fmt.Println("🌐 Servidor ejecutándose en http://localhost:8080")
	r.Run(":8080")
}
