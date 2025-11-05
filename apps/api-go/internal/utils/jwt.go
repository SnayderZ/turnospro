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
func GenerarToken(email string, rol string) (string, error) {
	claims := jwt.MapClaims{ // mapa de claims (contenido del token)
		"email": email,                                // quién es el usuario
		"rol":   rol,                                  // qué rol tiene (admin/operador)
		"exp":   time.Now().Add(2 * time.Hour).Unix(), // fecha de expiración (UNIX)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // token HS256
	return token.SignedString(secretKey())                     // firmar con clave secreta
}
