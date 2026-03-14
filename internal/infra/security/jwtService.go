package security

import (
	"crypto/rsa"
	"os"
	"plantao/internal/infra/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	expireTime int64
}

type JWTClaims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(cfg *config.Config) (*JWTService, error) {
	privateKeyBytes, err := os.ReadFile(cfg.JWT.PrivateKeyPath)

	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := os.ReadFile(cfg.JWT.PublicKeyPath)

	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyBytes))

	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyBytes))

	if err != nil {
		return nil, err
	}

	return &JWTService{
		privateKey: privateKey,
		publicKey:  publicKey,
		expireTime: cfg.JWT.ExpireTime,
	}, nil
}

func (j *JWTService) GenerateToken(userId string, role string) (string, error) {
	claims := JWTClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expireTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "intellisys",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(j.privateKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, jwt.ErrSignatureInvalid
		}

		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
