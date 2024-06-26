package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Hotkey 
var Hotkey = []byte("G0-store")

// SaltStr  
var SaltStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//Claims  
type Claims struct {
	UserID uint
	jwt.StandardClaims
}

// NextToken 
func NextToken(uid uint) string {
	fmt.Printf("传入JWT的id:%v/n", uid)
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), 
			IssuedAt:  time.Now().Unix(),
			Issuer:    "root",       
			Subject:   "user token", 
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	tokenString, err := token.SignedString(Hotkey)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

// ParseToken 
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Hotkey, nil
	})
	if err != nil {
		fmt.Println(" token parse err:", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
