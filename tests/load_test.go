package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	rate := vegeta.Rate{Freq: 10, Per: vegeta.Second} // 10 requests per second
	duration := 10 * time.Second                      // Run test for 10 seconds

	// Prepare the targeter
	targeter := vegeta.NewStaticTargeter(generateTarget())

	// Prepare the attacker
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "FileUploadTest") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Duration: %s\n", metrics.Duration)
	fmt.Printf("Success: %f\n", metrics.Success*100)
	fmt.Printf("Status Codes: %v\n", metrics.StatusCodes)
	fmt.Printf("Errors: %v\n", metrics.Errors)
}

func generateTarget() vegeta.Target {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := os.Open("../tests/test_files/file.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("files", filepath.Base(file.Name()))
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("Failed to copy file: %v", err)
	}
	writer.Close()

	return vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8989/upload",
		Header: map[string][]string{
			"Content-Type": {writer.FormDataContentType()},
		},
		Body: body.Bytes(),
	}
}
