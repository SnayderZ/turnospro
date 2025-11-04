package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"turnospro/api-go/internal/db"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// === 1️⃣ Cargar variables de entorno ===
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Println("⚠️  No se encontró .env, usando variables del sistema...")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ ERROR: no se encontró DB_URL")
	}

	// === 2️⃣ Conexión a PostgreSQL ===
	database, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("❌ Error al abrir BD:", err)
	}
	defer database.Close()

	err = database.Ping()
	if err != nil {
		log.Fatal("❌ Error al conectar con PostgreSQL:", err)
	}
	fmt.Println("✅ Conectado correctamente a PostgreSQL")

	// === 3️⃣ Crear instancia del repositorio SQLC ===
	queries := db.New(database)
	ctx := context.Background()

	// === 4️⃣ Insertar un usuario de prueba ===
	nuevo, err := queries.CreateUser(ctx, db.CreateUserParams{
		Nombre: "Usuario Demo",
		Email:  "demo@sqlc.com",
		Hash:   "clave_encriptada",
		Rol:    "admin",
	})
	if err != nil {
		log.Fatal("❌ Error al insertar usuario:", err)
	}

	fmt.Println("✅ Usuario insertado:", nuevo.Email)

	// === 5️⃣ Servidor básico de prueba ===
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Servidor activo 🚀",
		})
	})

	fmt.Println("🌐 Servidor ejecutándose en http://localhost:8080")
	r.Run(":8080")
}
