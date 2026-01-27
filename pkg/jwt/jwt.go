package jwt

import (
	"charonoms/internal/infrastructure/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	RoleID       uint   `json:"role_id"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID, roleID uint, username string, isSuperAdmin bool, cfg config.JWTConfig) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(cfg.ExpireHours) * time.Hour)

	claims := Claims{
		UserID:       userID,
		Username:     username,
		RoleID:       roleID,
		IsSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string, cfg config.JWTConfig) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token
func RefreshToken(tokenString string, cfg config.JWTConfig) (string, error) {
	claims, err := ParseToken(tokenString, cfg)
	if err != nil {
		return "", err
	}

	// 重新生成 Token
	return GenerateToken(claims.UserID, claims.RoleID, claims.Username, claims.IsSuperAdmin, cfg)
}
