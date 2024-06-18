package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	// home, err := os.UserHomeDir()
	// uploadPath := filepath.Join(home, "/.upload")
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %s\n", err)
		return
	}
	uploadPath := filepath.Join(home, ".uploads")

	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating upload directory:", err)
		return
	}
	// fmt.Println("Home directory:", uploadPath)
	// Serving static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Handle file uploads
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		fileUploadHandler(w, r, uploadPath)
	})
	fmt.Println("Server started on :8989")
	err = http.ListenAndServe(":8989", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request, uploadPath string) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the ParseMultipartForm method is called before retrieving files
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve files from the posted form-data
	files := r.MultipartForm.File["files"]
	if files == nil {
		http.Error(w, "No files received", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	var successfulUploads []string
	// Below loop ensures multi-file upload

	for _, fileHeader := range files {
		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()
			file, err := fileHeader.Open()
			now := time.Now()
			format := "2006-01-02 15:04:05"

			// Format the timestamp and print it

			if err != nil {
				fmt.Println("Error retrieving the file", err)
				// continue
			}
			defer file.Close()

			newPath := filepath.Join(uploadPath, filepath.Base(fileHeader.Filename))
			newFile, err := os.Create(newPath)
			if err != nil {
				fmt.Println("Error creating the file", err)
				return
			}
			defer newFile.Close()

			fileName := fileHeader.Header.Get("Content-Disposition")

			// Extract file name from Content-Disposition header (optional parsing)
			var actualFileName string
			if fileName != "" {
				// Basic parsing (can be improved for robustness)
				parts := strings.Split(fileName, `;`)
				for _, part := range parts {
					if strings.Contains(part, "filename=") {
						actualFileName = strings.TrimSpace(strings.SplitN(part, "=", 2)[1])
						break
					}
				}
			}
			timestamp := now.Format(format)
			fmt.Printf("%s: File being uploaded: %s\n\n", timestamp, actualFileName)

			bytesWritten, err := io.Copy(newFile, file)
			fmt.Fprintf(w, "File size is %d\n", bytesWritten)
			if err != nil {
				http.Error(w, "Error saving the file", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			successfulUploads = append(successfulUploads, fileHeader.Filename)
			logUploadDetails(uploadPath, filepath.Base(fileHeader.Filename), bytesWritten)
			fmt.Fprintf(w, "File uploaded successfully: %+v", fileHeader.Filename)

		}(fileHeader)
	}
	wg.Wait()

	// Report back to client
	if len(successfulUploads) > 0 {
		fmt.Fprintf(w, "Successfully uploaded files: %v", successfulUploads)
	} else {
		http.Error(w, "Failed to upload any files", http.StatusInternalServerError)
	}

	/*
		Below lines of code for single file upload
	*/
	// file, header, err = r.FormFile("file")
	// if err != nil {
	// 	http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
	// 	fmt.Println(err)
	// 	return
	// }
	// defer file.Close()
	// newPath := filepath.Join("/Users/harishkumarpillai/.uploads", filepath.Base(header.Filename))

	// fmt.Printf("Uploaded File: %+v\n", header.Filename)
	// fmt.Printf("File Size: %+v\n", header.Size)
	// fmt.Printf("MIME Header: %+v\n", header.Header)

	// newFile, err := os.Create(newPath)
	// if err != nil {
	// 	http.Error(w, "Error creating the file", http.StatusInternalServerError)
	// 	fmt.Println(err)
	// 	return
	// }
	// defer newFile.Close()

	// bytesWritten, err := io.Copy(newFile, file)
	// if err != nil {
	// 	http.Error(w, "Error saving the file", http.StatusInternalServerError)
	// 	fmt.Println(err)
	// 	return
	// }
	// logUploadDetails(header.Filename, bytesWritten)
	// fmt.Fprintf(w, "File uploaded successfully: %+v", header.Filename)
}

func logUploadDetails(filePath string, newPath string, fileSize int64) {
	logFile, err := os.OpenFile(filepath.Join(filePath, "upload_logs.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	now := time.Now()
	logEntry := fmt.Sprintf("Upload: %s, Size: %d bytes, Date: %s\n", newPath, fileSize, now.Format("2006-01-02 15:04:05"))
	if _, err = logFile.WriteString(logEntry); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
