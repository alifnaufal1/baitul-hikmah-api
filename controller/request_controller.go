package controller

import (
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("my_scret_key")

func RequestController(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization") 
	if authHeader == "" {
		utils.HandleAnyError("no token provided", w, 401)
		return
	}
	
	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) != 2 || tokenString[0] != "Bearer" {
		utils.HandleAnyError("invalid token format", w, 401)
		return
	}
	
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SECRET_KEY, nil
	})

	if err != nil || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			utils.HandleAnyError("expired token", w, 401)
			return
		}
		utils.HandleAnyError("invalid token", w, 401)
		return
	}
	
	var userID int
	switch v := claims["user_id"].(type) {
	case float64:
		userID = int(v)
	case string:
		id, err := strconv.Atoi(v)
		if err == nil {
			userID = id
		}
	}
	
	if repo.IsBlacklistToken(token.Raw, userID) {
		utils.HandleAnyError("blacklist token", w, 401)
		return
	}

	validResponse := types.APIResponse{
		Code: 200,
		Results: types.ValidTokenResult{
			Message: "valid token",
		},
		Status: "success",
	}
	
	utils.JSONTemplate(w, 200, validResponse)
}