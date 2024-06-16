# Go File Upload Server

This repository contains a simple but powerful implementation of a file upload server written in Go. It demonstrates how to handle multiple file uploads using `multipart/form-data` format and how to save these files securely on the server filesystem. This project is platform-independent and can run on macOS, Windows, and Linux.


## Features

- Multi-file upload handling using HTTP POST requests.
- Concurrent file processing to optimize performance.
- Cross-platform path handling for file operations.
- Detailed logging of file upload details and errors.

Run the file upload server 
go run main.go

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have Go installed on your machine. To install Go, follow the instructions on the [official Go website](https://golang.org/dl/).

### Installing

1. **Clone the repository**

   ```bash
   git clone https://github.com/pillaiharish/go-file-upload-server.git
   ```

## Running load tests
```bash
cd tests
go run load_test.go
```

## Running vegeta
```bash
cd tests/vegeta
./attack.sh
```


## Install vegeta
```bash
go install github.com/tsenart/vegeta@latest
```

## Install Gnuplot:
- On macOS: brew install gnuplot
- On Linux: sudo apt-get install gnuplot
- On Windows: Download and install from the official website.


## Run the Attack and Plot the Data:
```bash
chmod +x tests/vegeta/attack.sh
./tests/vegeta/attack.sh &
gnuplot -persist tests/vegeta/plot.gp
```

## Cleanup:
After you're done, you might want to remove the named pip
```bash
rm /tmp/results.bin
```
