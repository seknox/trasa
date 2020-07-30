package my

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/seknox/trasa/core/redis"
	"github.com/seknox/trasa/models"
	"github.com/seknox/trasa/utils"
	"github.com/sirupsen/logrus"
)

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 5000000000)
	err := r.ParseForm()
	if err != nil {
		// redirect or set error status code.
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "File too large", "File not uploaded", nil)
		return
	}

	userContext := r.Context().Value("user").(models.UserContext)
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		logrus.Errorf("error retrieving the file: %v", err)
		utils.TrasaResponse(w, 200, "failed", "no file uploaded", "File not uploaded", nil)
		return
	}
	defer file.Close()

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	err = os.MkdirAll(filepath.Join("/tmp/trasa/trasagw/shared/", userContext.User.ID), os.ModePerm)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not create shared directory ", "File not uploaded", nil)
		return
	}

	f, err := os.Create(filepath.Join("/tmp/trasa/trasagw/shared/", userContext.User.ID, handler.Filename))
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not create a file", "File not uploaded", nil)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	utils.TrasaResponse(w, 200, "success", "file uploaded", fmt.Sprintf(`File %s uploaded`, handler.Filename), nil)
}

func GetDownloadableFileList(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	userDir := filepath.Join("/tmp/trasa/trasagw/shared/", userContext.User.ID)
	fileList, err := ioutil.ReadDir(userDir)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not get file list", "GetDownloadableFileList")
		return
	}
	fileNameList := []string{}
	for _, fileI := range fileList {
		if !fileI.IsDir() {
			fileNameList = append(fileNameList, fileI.Name())
		}

	}
	utils.TrasaResponse(w, 200, "success", "dir list", "GetDownloadableFileList", fileNameList)
}

func GetFileDownloadToken(w http.ResponseWriter, r *http.Request) {
	userContext := r.Context().Value("user").(models.UserContext)
	userID := userContext.User.ID
	authKey := utils.GetRandomID(17)
	err := redis.Store.Set(authKey, time.Second*600, "userID", userID)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not save token", "file not uploaded")
		return
	}
	w.Header().Set("sskey", authKey)
	utils.TrasaResponse(w, 200, "success", "", fmt.Sprintf(`File downloaded`), nil)
	return

}

func FileDownloadHandler(w http.ResponseWriter, r *http.Request) {

	fileNameEnc := chi.URLParam(r, "fileName")
	fileName, err := url.PathUnescape(fileNameEnc)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "bad file name", "FileDownloadHandler", nil)
		return
	}
	ssKey := chi.URLParam(r, "sskey")
	userID, err := redis.Store.Get(ssKey, "userID")

	if err != nil || ssKey == "" || userID == "" {
		utils.TrasaResponse(w, 200, "failed", "invalid sskey", "FileDownloadHandler", nil, nil)
		return
	}

	if strings.Contains(fileName, "/") {
		utils.TrasaResponse(w, 200, "failed", "invalid filename", "FileDownloadHandler", nil, nil)
		return
	}

	filePath := filepath.Join("/tmp/trasa/trasagw/shared/", userID, fileName)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	http.ServeFile(w, r, filePath)
}
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//http.StripPrefix("/api/v1/download_file", http.FileServer(http.Dir("../trasagw/shared"))).ServeHTTP
	userContext := r.Context().Value("user").(models.UserContext)
	userID := userContext.User.ID
	fileNameEnc := chi.URLParam(r, "fileName")
	fileName, err := url.PathUnescape(fileNameEnc)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "bad file name", "file delete", nil)
		return
	}

	if strings.Contains(fileName, "/") {
		utils.TrasaResponse(w, 200, "failed", "invalid filename", "FileDownloadHandler")
		return
	}

	filePath := filepath.Join("/tmp/trasa/trasagw/shared/", userID, fileName)
	err = os.Remove(filePath)
	if err != nil {
		logrus.Error(err)
		utils.TrasaResponse(w, 200, "failed", "could not delete file", "file delete")
		return
	}
	utils.TrasaResponse(w, 200, "success", "file deleted", "file delete")

}
