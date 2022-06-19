package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Token_id     string `json:"token_id"`
	User_id      string `json:"user_id"`
	Auth_token   string `json:"auth_token"`
	Generated_at string `json:"generated_at"`
	Expires_at   string `json:"expires_at"`
}

func (db *DB) GenerateToken(phone_number string, password string) (map[string]interface{}, error) {

	userInfo, err := db.GetUserByPhoneNumber(phone_number)

	if err != nil {
		return nil, err
	}

	userId := userInfo.User_id
	accountPassword := userInfo.Password

	log.Println(accountPassword)
	log.Println(password)

	err = bcrypt.CompareHashAndPassword([]byte(accountPassword), []byte(password))

	if err != nil {
		return nil, errors.New("invalid password")
	}

	queryString := "insert into authentication_tokens(user_id, auth_token, generated_at, expires_at) values (?, ?, ?, ?)"
	stmt, err := db.Prepare(queryString)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	randomToken := make([]byte, 32)

	_, err = rand.Read(randomToken)

	if err != nil {
		return nil, err
	}

	authToken := base64.URLEncoding.EncodeToString(randomToken)

	const timeLayout = "2006-01-02 15:04:05"

	dt := time.Now()
	expirtyTime := time.Now().Add(time.Minute * 60)

	generatedAt := dt.Format(timeLayout)
	expiresAt := expirtyTime.Format(timeLayout)

	_, err = stmt.Exec(userId, authToken, generatedAt, expiresAt)

	if err != nil {
		return nil, err
	}

	tokenDetails := map[string]interface{}{
		"token_type":   "Bearer",
		"auth_token":   authToken,
		"generated_at": generatedAt,
		"expires_at":   expiresAt,
	}

	return tokenDetails, nil

}

func (db *DB) ValidateToken(authToken string) (map[string]interface{}, error) {

	userTokenInfo, err := db.GetUserByAuthToken(authToken)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid access token.\r\n")
		}

		return nil, err
	}

	const timeLayout = "2006-01-02 15:04:05"

	expiryTime, _ := time.Parse(timeLayout, userTokenInfo.Expires_at)
	currentTime, _ := time.Parse(timeLayout, time.Now().Format(timeLayout))

	if expiryTime.Before(currentTime) {
		return nil, errors.New("The token is expired.\r\n")
	}

	userDetails := map[string]interface{}{
		"user_id":      userTokenInfo.User.User_id,
		"first_name":   userTokenInfo.User.First_name,
		"generated_at": userTokenInfo.Token.Generated_at,
		"expires_at":   userTokenInfo.Token.Expires_at,
	}

	return userDetails, nil
}
