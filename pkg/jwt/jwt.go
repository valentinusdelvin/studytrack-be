package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/valentinusdelvin/studytrack-be/internal/infra/env"
)

var ErrInvalidTokenType = errors.New("invalid token type")

type CustomJwtInterface interface {
	Create(
		id uint64,
		name string,
		email string,
	) (string, error)
	Decode(tokenString string, claims *Claims) error
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64 `json:"-"`
	Name   string `json:"nama"`
	Email  string `json:"email"`
}

type CustomJwtStruct struct {
	SecretKey   string
	ExpiredTime time.Duration
}

var Jwt = getJwt()

func getJwt() CustomJwtInterface {
	return &CustomJwtStruct{
		SecretKey:   env.AppEnv.JwtSecretKey,
		ExpiredTime: env.AppEnv.JwtExpTime,
	}
}

func (j *CustomJwtStruct) generateToken(
	id uint64,
	name string,
	email string,
	duration time.Duration,
) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "studytrack-be",
			Audience:  jwt.ClaimStrings{"studytrack-be"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
			Subject:   strconv.FormatUint(id, 10),
		},
		UserID: id,
		Name:   name,
		Email:  email,
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j *CustomJwtStruct) Create(
	id uint64,
	name string,
	email string,
) (string, error) {
	return j.generateToken(id, name, email, j.ExpiredTime)
}

func (j *CustomJwtStruct) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	},
		jwt.WithIssuer("studytrack-be"),
		jwt.WithAudience("studytrack-be"),
	)

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}
	claims.UserID, err = strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user id claim: %w", err)
	}

	return nil
}
