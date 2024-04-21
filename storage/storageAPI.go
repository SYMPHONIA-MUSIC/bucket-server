package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadMediaHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("UploadMediaHandler started")

	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		log.Printf("UploadMediaHandler error: unsupported method %v", r.Method)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		log.Printf("UploadMediaHandler error: invalid file - %v", err)
		return
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.WriteString(hash, time.Now().Format(time.RFC3339Nano)); err != nil {
		http.Error(w, "Failed to hash time", http.StatusInternalServerError)
		log.Printf("UploadMediaHandler error: failed to hash time - %v", err)
		return
	}
	if _, err := io.Copy(hash, file); err != nil {
		http.Error(w, "Failed to hash file", http.StatusInternalServerError)
		log.Printf("UploadMediaHandler error: failed to hash file - %v", err)
		return
	}

	hashedFilename := hex.EncodeToString(hash.Sum(nil)) + filepath.Ext(header.Filename)
	filepath := filepath.Join("data", hashedFilename)

	out, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Printf("UploadMediaHandler error: failed to save file - %v", err)
		return
	}
	defer out.Close()

	file.Seek(0, 0)
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		log.Printf("UploadMediaHandler error: failed to write file - %v", err)
		return
	}

	response := map[string]string{"hashData": hashedFilename}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("File uploaded successfully: %s in %v", hashedFilename, time.Since(start))
}

func FetchMediaHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("FetchMediaHandler started")

	if r.Method != "GET" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		log.Printf("FetchMediaHandler error: unsupported method %v", r.Method)
		return
	}

	fileHash := r.URL.Query().Get("hash")
	if fileHash == "" {
		http.Error(w, "File hash is required", http.StatusBadRequest)
		log.Printf("FetchMediaHandler error: file hash is required")
		return
	}

	filepath := filepath.Join("data", fileHash)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("FetchMediaHandler error: file not found - %v", fileHash)
		return
	}

	log.Printf("File fetched successfully: %s in %v", fileHash, time.Since(start))
	http.ServeFile(w, r, filepath)
}
