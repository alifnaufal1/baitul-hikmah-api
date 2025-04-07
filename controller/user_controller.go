package controller

import (
	"blog-api/middleware"
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func UserRegisterController(w http.ResponseWriter, r *http.Request) {
	var user types.User

	if err := utils.DecodeFromRequest(r.Body, &user, w); err != nil {return}
	
	if len(user.Password) > 8 {
		utils.HandleAnyError("password is too long", w, http.StatusBadRequest)
		return
	}
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	
	
	registeredUser, err := repo.CreateUser(user.Username, string(hashedPassword))
	if err != nil {
		utils.HandleAnyError("error registration -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}
	
	utils.SuccessResponse(w, 201, registeredUser, "success register")
}

func UserLoginController(w http.ResponseWriter, r *http.Request)  {
	var user types.User
	
	if err := utils.DecodeFromRequest(r.Body, &user, w); err != nil {return}
	
	dbUser, err := repo.GetUserByUsername(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleDataNotFound("this user not found", w)
			return
		}
		utils.HandleAnyError("error get post -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) 
	if err != nil {
		utils.HandleAnyError("invalid credentials", w, http.StatusUnauthorized)
		return
	}

	
	generatedToken, err := middleware.GenerateToken(dbUser.ID, dbUser.Role)
	if err != nil {
		utils.HandleAnyError("error generating token -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}

	loginResponse := types.APIResponse{
		Code: 200,
		Results: types.LoginResult{
			Message: "success login",
			Token: generatedToken,	
		},
		Status: "succes",
	}
	
	utils.JSONTemplate(w, 200, loginResponse)
}