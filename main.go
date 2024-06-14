package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	os.MkdirAll("/Users/harishkumarpillai/.uploads", os.ModePerm)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/upload", fileUploadHandler)
	fmt.Println("Server started on :8189")
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Failed to create multipart reader: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	var successfulUploads []string

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to read multipart data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if part.FileName() == "" {
			continue
		}

		wg.Add(1)
		go func(part *multipart.Part) {
			defer wg.Done()

			now := time.Now()
			format := "2006-01-02 15:04:05"
			timestamp := now.Format(format)
			fmt.Printf("%s: File being uploaded: %s\n", timestamp, part.FileName())

			newPath := filepath.Join("/Users/harishkumarpillai/.uploads", filepath.Base(part.FileName()))
			newFile, err := os.Create(newPath)
			if err != nil {
				fmt.Println("Error creating the file", err)
				return
			}
			defer newFile.Close()

			bytesWritten, err := io.Copy(newFile, part)
			if err != nil {
				fmt.Println("Error saving the file", err)
				return
			}

			successfulUploads = append(successfulUploads, part.FileName())
			logUploadDetails(part.FileName(), bytesWritten)
			fmt.Fprintf(w, "File uploaded successfully: %+v\n", part.FileName())
		}(part)
	}

	wg.Wait()

	if len(successfulUploads) > 0 {
		fmt.Fprintf(w, "Successfully uploaded files: %v", successfulUploads)
	} else {
		http.Error(w, "Failed to upload any files", http.StatusInternalServerError)
	}
}

func logUploadDetails(filename string, fileSize int64) {
	logFile, err := os.OpenFile("/Users/harishkumarpillai/.uploads/upload_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	now := time.Now()
	logEntry := fmt.Sprintf("Upload: %s, Size: %d bytes, Date: %s\n", filename, fileSize, now.Format("2006-01-02 15:04:05"))
	if _, err = logFile.WriteString(logEntry); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
