package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"rajdeepm.xyz/ofm/models"
)

// Serve a file or return an error
func ServeFile(c *gin.Context, path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			jsonError(c, http.StatusNotFound, "File not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	if fileInfo.IsDir() {
		jsonError(c, http.StatusTeapot, "Cannot fetch a directory as a file")
		return
	}

	c.File(path) // Use Gin's File method to serve the file
}

// getFile handles GET requests for files
func GetFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "config.toml" {
		jsonError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	path := filepath.Join(".", filename)
	ServeFile(c, path)
}

// listDirectory lists files in the requested directory
func ListDirectory(c *gin.Context, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to read directory")
		return
	}

	var metaList []models.Meta
	for _, entry := range entries {
		if entry.Name() == "config.toml" {
			continue
		}

		absPath, _ := filepath.Abs(filepath.Join(dir, entry.Name()))
		info, _ := entry.Info()
		meta := models.Meta{
			Filename: entry.Name(),
			Absolute: absPath,
			IsDir:    entry.IsDir(),
			Size:     info.Size(),
		}
		metaList = append(metaList, meta)
	}

	c.JSON(http.StatusOK, metaList) // Send the JSON response
}

// listing handles directory listing requests
func Listing(c *gin.Context) {
	dir := c.Param("dir")
	if dir == "$pwd" {
		dir = "."
	}

	fileInfo, err := os.Stat(dir)
	if err != nil || !fileInfo.IsDir() {
		jsonError(c, http.StatusNotFound, "Directory not found")
		return
	}

	ListDirectory(c, dir)
}

// uploadFile handles file uploads
func UploadFile(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "config.toml" {
		jsonError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	path := filepath.Join(".", filename)

	// Create all necessary parent directories
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to create directories")
		return
	}

	file, err := os.Create(path)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer file.Close()

	if _, err = io.Copy(file, c.Request.Body); err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to write to file")
		return
	}

	c.JSON(http.StatusOK, ErrorResponse{Status: http.StatusOK, Message: "File uploaded successfully"})
}

// deleteFile handles file deletions
func DeleteFile(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "config.toml" {
		jsonError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	path := filepath.Join(".", filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		jsonError(c, http.StatusNotFound, "File or directory not found")
		return
	}

	if err := os.RemoveAll(path); err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to delete file or directory")
		return
	}

	c.JSON(http.StatusOK, ErrorResponse{Status: http.StatusOK, Message: "File or directory deleted successfully"})
}

// renameFile handles file renaming using a JSON request
func RenameFile(c *gin.Context) {
	oldName := c.Param("oldname")
	var request struct {
		NewName string `json:"new"`
	}

	if err := c.BindJSON(&request); err != nil {
		jsonError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	newName := request.NewName
	if newName == "" {
		jsonError(c, http.StatusBadRequest, "New filename is required")
		return
	}

	if oldName == "config.toml" || newName == "config.toml" {
		jsonError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	oldPath := filepath.Join(".", oldName)
	newPath := filepath.Join(".", newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		jsonError(c, http.StatusNotFound, "File not found")
		return
	}

	if err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm); err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to create directories for new path")
		return
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to rename file")
		return
	}

	c.JSON(http.StatusOK, ErrorResponse{Status: http.StatusOK, Message: "File renamed successfully"})
}
