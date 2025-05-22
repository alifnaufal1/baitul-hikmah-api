package controller

import (
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
)

func getProfileImageName(w http.ResponseWriter, r *http.Request, uploadedFilename string) (string, error) {
	userIDCtx := r.Context().Value(types.UserKey)
	if userIDCtx == nil {
		utils.HandleAnyError("Unauthorized", w, http.StatusUnauthorized)
		return "", nil
	}

	userID, ok := userIDCtx.(int)
	if !ok {
		utils.HandleAnyError("Invalid user ID", w, http.StatusInternalServerError)
		return "", nil
	}

	registeredUser, err := repo.GetUserById(userID)
	if err != nil {
		utils.HandleDataNotFound("Cannot fetch user", w)
		return "", err
	}

	fileName := fmt.Sprintf("%d-%s%s", userID, registeredUser.Username, filepath.Ext(uploadedFilename))
	
	return fileName, nil
}

func getPostImageName(postID int, postTitle string, ext string) (string) {
	postTitle = strings.ToLower(postTitle)
	reg, _ := regexp.Compile(`[^a-z0-9\s]+`)
	postTitle = reg.ReplaceAllString(postTitle, " ")
	if strings.Contains(postTitle, " ") {
		postTitle = strings.ReplaceAll(postTitle, " ", "_")
	}
	return fmt.Sprintf("%d-%s%s", postID, postTitle, ext)	
}

func UploadImageController(w http.ResponseWriter, r *http.Request, dstDir types.DirName) (string, error) {
	err := r.ParseMultipartForm(1024)
	if err != nil {return "", err}
	
	file, header, err := r.FormFile(dstDir.ImageType + "-image")
	if file == nil {return "", errors.New("file is required")}
	if err != nil {return "", err}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	resizedImg := imaging.Resize(img, 800, 0, imaging.Lanczos)
	
	dir, err := os.Getwd()
	if err != nil {return "", err}
	
	var fileName string
	if dstDir.ImageType == "profile" {
		fileName, err = getProfileImageName(w, r, header.Filename)
		if err != nil || fileName == "" {
			utils.HandleDataNotFound("cannot fetch user", w)
			return "", err
		}	
	} else if (dstDir.ImageType == "post") {
		fileName = getPostImageName(dstDir.IdFileName, dstDir.ImageName, filepath.Ext(header.Filename))
	} else {return "", err}
	
	fileLocation := filepath.Join(dir, "uploads/" + dstDir.ImageType, fileName)
	
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {return "", err}
	defer targetFile.Close()
	
	switch filepath.Ext(header.Filename) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(targetFile, resizedImg, &jpeg.Options{Quality: 85})
	case ".png":
		err = png.Encode(targetFile, resizedImg)
	default:
		return "", errors.New("unsupported image format (only receive .jpg, .jpeg, & .png)")
	}
	if err != nil {return "", err}	

	dbFileName := fmt.Sprintf("%s/uploads/%s/%s", utils.Domain, dstDir.ImageType, fileName)

	return dbFileName, nil
}