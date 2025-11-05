package utils

import (
	"os"   // para leer variables de entorno (JWT_SECRET)
	"time" // para manejar expiraciones (exp)

	"github.com/golang-jwt/jwt/v5" // librería JWT
)

// secretKey obtiene la clave desde el entorno; si no existe, usa un valor por defecto.
// En producción SIEMPRE debes poner JWT_SECRET en el .env / variables de entorno.
func secretKey() []byte {
	if s := os.Getenv("../../../../../.env"); s != "" { // si existe JWT_SECRET en el entorno...
		return []byte(s)
	}
	return []byte("super_secret_key")
}

// GenerarToken crea un JWT HS256 que expira en 2 horas.
// Guardamos claims mínimos: email y rol (útiles para autorización).
func GenerarToken(id int32, email, rol string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"user_id": id, // ✅ agregamos el id del usuario
		"email":   email,
		"rol":     rol,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
