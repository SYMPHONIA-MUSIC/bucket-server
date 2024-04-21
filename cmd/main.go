package main

import (
	"bucket-server/authentication"
	"bucket-server/storage"
	"bucket-server/util"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	cfg, err := util.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fileStorageDir := "data"

	if _, err := os.Stat(fileStorageDir); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist, attempting to create it", fileStorageDir)
		err := os.MkdirAll(fileStorageDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", fileStorageDir, err)
		}
		log.Printf("Directory %s created successfully", fileStorageDir)
	}

	http.HandleFunc("/upload", logRequest(authentication.ValidateKey(cfg.APIKey, storage.UploadMediaHandler)))
	http.HandleFunc("/fetch", logRequest(authentication.ValidateKey(cfg.APIKey, storage.FetchMediaHandler)))

	if cfg.UseHTTPS {
		log.Printf("Starting HTTPS server on port 443")
		log.Fatal(http.ListenAndServeTLS(":443", cfg.HTTPSCertPath, cfg.HTTPSKeyPath, nil))
	} else {
		log.Printf("Starting HTTP server on port " + cfg.ServerPort)
		log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
	}
}

func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler(w, r)
		log.Printf("[%s] %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	}
}
