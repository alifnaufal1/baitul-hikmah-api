package controller

import (
	"blog-api/middleware"
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

func UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		utils.HandleAnyError("error parse multipart form -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}
	
	uploadedFile, headerFile, err := r.FormFile("profile-image")
	if err != nil {
		utils.HandleAnyError("file required", w, 402)
		return
	}
	
	defer uploadedFile.Close()
	
	dir, err := os.Getwd()
	if err != nil {
		utils.HandleAnyError("error get directory -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}

	userIDCtx := r.Context().Value(types.UserKey)
  if userIDCtx == nil {
    utils.HandleAnyError("Unauthorized", w, http.StatusUnauthorized)
    return
  }

  userID, ok := userIDCtx.(int)
  if !ok {
		utils.HandleAnyError("Invalid user ID", w, http.StatusInternalServerError)
		return
	}

	fmt.Println("user_id from context:", userID)
	
	registeredUser, err := repo.GetUsernameById(userID)
	if err != nil {
		utils.HandleDataNotFound("Cannot fetch user", w)
		return
	}
	
	fileName := fmt.Sprintf("%d-%s%s", userID, registeredUser.Username, filepath.Ext(headerFile.Filename))
	
	fileLocation := filepath.Join(dir, "uploads/profile", fileName)

	_, err = repo.AddProfileImage(userID, "http://localhost:90/uploads/profile/" + fileName)
	if err != nil {
		utils.HandleAnyError("cannot save profile image to db -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}
	
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		utils.HandleAnyError("cannot save profile image -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}
	
	defer targetFile.Close()
	
	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		utils.HandleAnyError("cannot copy profile image -> " + err.Error(), w, http.StatusInternalServerError)
		return
	}

	uploadResponse := types.APIResponse{
		Code: 200,
		Results: "success",
		Status: "success",
	}
	
	utils.JSONTemplate(w, 200, uploadResponse)
}