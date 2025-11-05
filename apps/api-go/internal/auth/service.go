package auth

import (
	"context"                         // contexto para cancelaciones/timeouts
	"database/sql"                    // para recibir *sql.DB
	"errors"                          // para devolver errores legibles
	"turnospro/api-go/internal/db"    // código generado por sqlc (Queries, modelos)
	"turnospro/api-go/internal/utils" // utilidades (JWT)

	"golang.org/x/crypto/bcrypt" // hashing/verificación de contraseñas
)

// AuthService orquesta operaciones de login/registro contra la BD
type AuthService struct {
	q *db.Queries // repositorio tipado generado por sqlc
}

// NewAuthService crea una instancia inyectando la conexión
func NewAuthService(dbConn *sql.DB) *AuthService {
	return &AuthService{q: db.New(dbConn)} // db.New(*sql.DB) -> *Queries
}

// Register crea un usuario nuevo encriptando la contraseña con bcrypt
func (s *AuthService) Register(ctx context.Context, nombre, email, password, rol string) (db.Usuario, error) {
	// bcrypt.GenerateFromPassword recibe []byte y un "cost" (DefaultCost = seguro por defecto)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.Usuario{}, err // si falla hashing, devuelve error
	}
	// Insertar el usuario usando la consulta tipada de sqlc
	return s.q.CreateUser(ctx, db.CreateUserParams{
		Nombre: nombre,
		Email:  email,
		Hash:   string(hashed), // guardamos el HASH, jamás el texto plano
		Rol:    rol,
	})
}

// Login valida credenciales y devuelve un JWT si son correctas
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// 1) Buscar el usuario por email
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("usuario no encontrado")
	}
	// 2) Comparar contraseña en texto con el hash almacenado
	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return "", errors.New("contraseña incorrecta")
	}
	// 3) Generar token JWT con email y rol
	token, err := utils.GenerarToken(int32(user.ID), user.Email, user.Rol)
	if err != nil {
		return "", err
	}
	return token, nil
}
