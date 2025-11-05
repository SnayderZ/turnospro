// apps/api-go/internal/auth/handler.go
package auth // mismo paquete

import (
	"context"      // para pasar context al servicio
	"database/sql" // para el *sql.DB
	"net/http"     // códigos HTTP

	"github.com/gin-gonic/gin" // framework HTTP
)

// RegisterRoutes registra las rutas /auth/* en el router principal
func RegisterRoutes(r *gin.Engine, dbConn *sql.DB) {
	service := NewAuthService(dbConn) // instancia servicio con repos sqlc

	// POST /auth/register — crea usuario
	r.POST("/auth/register", func(c *gin.Context) {
		var req struct { // DTO sencillo para deserializar JSON
			Nombre   string `json:"nombre"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil { // valida JSON
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
		user, err := service.Register(context.Background(), req.Nombre, req.Email, req.Password, "operador")
		if err != nil { // puede fallar por email duplicado u otro error de BD
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Usuario registrado correctamente",
			"user":    user.Email,
		})
	})

	// POST /auth/login — devuelve token
	r.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
		token, err := service.Login(context.Background(), req.Email, req.Password)
		if err != nil { // usuario no encontrado o contraseña incorrecta
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
