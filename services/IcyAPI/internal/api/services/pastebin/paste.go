package pastebin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (c *Controller) CreatePaste(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content" validate:"required"`
		Syntax   string `json:"syntax"`
		ImageURL string `json:"image_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := c.Validator.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	paste := Paste{
		Title:       c.Sanitizer.Sanitize(req.Title),
		Content:     c.Sanitizer.Sanitize(req.Content),
		Syntax:      req.Syntax,
		ImageURL:    req.ImageURL,
		DeleteToken: uuid.NewString(),
	}
	

	if err := c.PostgresClient.DB.Create(&paste).Error; err != nil {
		http.Error(w, "Failed to save paste", http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{
		"message":      "Paste created",
		"id":           paste.ID,
		"delete_token": paste.DeleteToken,
	})
	
}

func (c *Controller) GetPaste(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/paste/"):]
	var paste Paste
	if err := c.PostgresClient.DB.First(&paste, "id = ?", id).Error; err != nil {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	paste.Views++
	c.PostgresClient.DB.Save(&paste)

	if r.URL.Query().Get("raw") == "true" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(paste.Content))
		return
	}

	JsonResponse(w, paste)
}

func (c *Controller) UpdatePaste(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/paste/update/"):]
	var paste Paste
	if err := c.PostgresClient.DB.First(&paste, "id = ?", id).Error; err != nil {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Syntax   string `json:"syntax"`
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.Title != "" {
		paste.Title = c.Sanitizer.Sanitize(req.Title)
	}
	if req.Content != "" {
		paste.Content = c.Sanitizer.Sanitize(req.Content)
	}
	if req.Syntax != "" {
		paste.Syntax = req.Syntax
	}
	if req.ImageURL != "" {
		paste.ImageURL = req.ImageURL
	}
	paste.UpdatedAt = time.Now()

	c.PostgresClient.DB.Save(&paste)
	JsonResponse(w, map[string]string{
		"message": "Paste updated",
	})
}

func (c *Controller) DeletePaste(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/paste/delete/"):]
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing delete token", http.StatusBadRequest)
		return
	}

	var paste Paste
	if err := c.PostgresClient.DB.First(&paste, "id = ?", id).Error; err != nil {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	if paste.DeleteToken != token {
		http.Error(w, "Invalid delete token", http.StatusUnauthorized)
		return
	}

	if err := c.PostgresClient.DB.Delete(&paste).Error; err != nil {
		http.Error(w, "Failed to delete paste", http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{
		"message": "Paste deleted",
	})
}

func (c *Controller) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Enforce maximum body size
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageURL, err := c.uploadImageToMinio(file, handler.Filename)
	if err != nil {
		http.Error(w, "failed to upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{
		"message": "Image uploaded",
		"url":     imageURL,
	})
}

func (c *Controller) uploadImageToMinio(file multipart.File, originalFilename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buf[:n])

	// Reset reader to beginning
	if seeker, ok := file.(io.Seeker); ok {
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("cannot seek uploaded file")
	}

	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file is not an image, detected: %s", contentType)
	}

	ext := filepath.Ext(originalFilename)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			ext = ".bin"
		}
	}

	newFilename := uuid.New().String() + ext

	_, err = c.MinoClient.Client.PutObject(ctx, "pastebin", newFilename, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:9000/%s/%s", "pastebin", newFilename), nil
}
