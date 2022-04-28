package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

const JWT_SECRET = "AKSAN-4RT"

type JwtStruct struct {
	CashierId float64
	Name      string // unique device mime
	Expired   interface{}
}

func CreateToken(id, name string) (token string, err error) {
	atClaims := jwt.MapClaims{}
	idx, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	atClaims["CashierId"] = idx
	atClaims["Name"] = name
	atClaims["Expired"] = time.Now().Add(time.Hour * 72).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err = at.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}
	return token, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(tokenString string) (result JwtStruct, err error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return result, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		result.CashierId = claims["CashierId"].(float64)
		result.Name = claims["Name"].(string)
		result.Expired = claims["Expired"]
	}
	var tm time.Time
	switch iat := result.Expired.(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}
	if (tm.Unix() - time.Now().Unix()) <= 0 {
		return result, errors.New("token expired, please login again")
	}
	return result, err
}
