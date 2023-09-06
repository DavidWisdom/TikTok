package controller
import (	
  "github.com/dgrijalva/jwt-go"
  "time"
  "fmt"
)
// 生成 Token
func GetToken(username, password string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"password": password,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with a secret key
	signedToken, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
// 解析 Token
func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("my_secret_key"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func GetInfo(token string) (string, string, error) {
	toke, err := ParseToken(token)
	if err != nil {
		return "", "", err
	}
	claims, ok := toke.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("Unexpected signing method: %v", toke.Header["alg"])
	}
	username := claims["username"].(string)
	password := claims["password"].(string)
	return username, password, nil
}