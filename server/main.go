package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const DEFAULT_SERVER_PORT = "8965"
const SERVER_FILE = "./server.pid"

func main() {
	StartServer(DEFAULT_SERVER_PORT, SERVER_FILE)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Query().Get("path")
	rangeHeader := r.Header.Get("Range")

	var start, end int64

	if filepath == "" {
		http.Error(w, "filepath param misssing!", http.StatusBadRequest)
		return
	}

	if isValidPath(filepath) {
		file, err := os.Open(filepath)

		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to open file!", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()

		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to read file info!", http.StatusInternalServerError)
			return
		}

		fileSize := fileInfo.Size()

		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
		w.Header().Set("Content-Type", "application/octet-stream")

		/*
			Details regarding "Range" Header (used to request partial content for pause-resume support in downloads)

			According to the HTTP/1.1 specification, RFC 2616, section 14.35.1,
			the Range header specifies a byte range as bytes-unit = "bytes" "=" byte-range,
			where byte-range is defined as first-byte-pos "-" [last-byte-pos].

			Here, first-byte-pos and last-byte-pos are both inclusive, and last-byte-pos
			is optional. If last-byte-pos is not provided, it means the range extends to the end of the file.
		*/
		if rangeHeader != "" {
			_, err := fmt.Scanf(rangeHeader, "bytes=%d-%d", &start, &end)
			if err != nil {
				_, err := fmt.Scanf(rangeHeader, "bytes=%d-", &start)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				end = fileSize - 1
			}

			if start > end {
				http.Error(w, "Invalid Range", http.StatusBadRequest)
			}
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
			w.WriteHeader(http.StatusPartialContent)
			file.Seek(start, io.SeekStart)

			/*
				By default, io.Copy uses a buffer of size 32KB. It reads up to 32KB from the file, writes it to the response writer, and then repeats the process until the entire file has been sent.
			*/

			io.Copy(w, io.LimitReader(file, end-start+1))
		} else {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
			io.Copy(w, file)
		}
	} else {
		http.Error(w, "Invalid filepath", http.StatusBadRequest)
		return
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pong!")
}

func StartServer(port string, serverFilePath string) {
	log.Print("Starting server...")

	err := os.WriteFile(serverFilePath, []byte(port), 0644)

	if err != nil {
		log.Fatalf("Unable to write server details: %s\n", err)
	}

	log.Printf("Server started on port: %s\n", port)

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/", handleFile)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func isValidPath(filePath string) bool {
	stat, err := os.Stat(filePath)

	if err != nil {
		return false
	}

	log.Println("File Stats: ")
	log.Printf("File Name: %s\n", stat.Name())
	log.Printf("File Size: %d KB\n\n", stat.Size()/1024)

	return true
}
