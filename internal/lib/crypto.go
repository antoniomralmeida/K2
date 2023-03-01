package lib

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const SignedKey = "uzUWld6Y0Ad6yUF8GU2gJGg8Q4wZaNNv"

func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder
	var (
		lowerCharSet   = "abcdedfghijklmnopqrstuvwxyz"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&*"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)
	time.Sleep(time.Microsecond)
	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func CreateJWTToken(id primitive.ObjectID, name string) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id.Hex()
	claims["name"] = name
	claims["exp"] = exp
	t, err := token.SignedString([]byte(SignedKey))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func DecodeToken(jwtcookie string) map[string]interface{} {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(jwtcookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SignedKey), nil
	})
	if err != nil {
		return nil
	}
	return claims
}

func ValidateToken(jwtcookie string) (err error) {

	if jwtcookie == "" {
		return InvalidTokenError
	}

	token, err := jwt.Parse(jwtcookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return []byte(SignedKey), nil
	})
	if err != nil {
		return
	}

	if token == nil {
		return InvalidTokenError
	}

	return nil
}
