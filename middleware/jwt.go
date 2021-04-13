package middleware

import (
	"errors"
	"github.com/biningo/boil-gin/global"
	"github.com/gin-gonic/gin"
)
import "github.com/dgrijalva/jwt-go"

/**
*@Author lyer
*@Date 4/13/21 15:52
*@Describe
**/

var (
	TokenInvalid = errors.New("invalid token")
	TokenExpired = errors.New("Token is expired")
	TokenError   = errors.New("Token error")
)

type Jwt struct {
	JwtSecret []byte
}

func NewJwt() *Jwt {
	return &Jwt{JwtSecret: []byte(global.G_CONFIG.Jwt.Secret)}
}

type CustomClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

func (j *Jwt) CreateToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtSecret)
}

func (j *Jwt) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtSecret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			}
			return nil, TokenError
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenInvalid
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(403, gin.H{"msg": "token不存在"})
			c.Abort()
			return
		}
		j := NewJwt()
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"msg": err.Error()})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
