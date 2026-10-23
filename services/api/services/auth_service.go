package services

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"medic-api/models"
	"medic-api/repositories"
)

var (
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrInactiveUser       = errors.New("akun pengguna tidak aktif")
	ErrJWTSecretMissing   = errors.New("JWT_SECRET belum dikonfigurasi")
)

type AuthService struct {
	userRepository *repositories.UserRepository
	jwtSecret      []byte
	tokenDuration  time.Duration
}

type AuthClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(
	userRepository *repositories.UserRepository,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      []byte(os.Getenv("JWT_SECRET")),
		tokenDuration:  15 * time.Minute,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (*models.LoginResponse, error) {
	if len(s.jwtSecret) == 0 {
		return nil, ErrJWTSecretMissing
	}

	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(s.tokenDuration)

	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "medic-api",
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Jangan sertakan password hash dalam respons login.
	user.PasswordHash = ""

	return &models.LoginResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		User:        *user,
	}, nil
}
