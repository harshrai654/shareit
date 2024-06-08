package main

import (
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

const DEFAULT_SERVER_PORT = "8966"

var RUNTIME_DIR = getRuntimeDirectory()
var SERVER_FILE = filepath.Join(RUNTIME_DIR, "server.pid")
var UNIX_SOCKET_FILE = filepath.Join(RUNTIME_DIR, "server.sock")

type filePathDetails struct {
	Secret string
	Otp    string
}

type socketPayload struct {
	FileAuth filePathDetails
	FilePath string
}

type FileVault struct {
	mu      sync.RWMutex
	fileMap map[string]filePathDetails
}

func (fv *FileVault) Read(key string) (filePathDetails, bool) {
	fv.mu.RLock()
	defer fv.mu.RUnlock()

	value, ok := fv.fileMap[key]

	return value, ok
}

func (fv *FileVault) Write(key string, value filePathDetails) {
	fv.mu.Lock()
	defer fv.mu.Unlock()

	fv.fileMap[key] = value
}

// According to OS returns the runtime directory
func getRuntimeDirectory() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "shareit")
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "shareit")
	}

	return filepath.Join(os.Getenv("HOME"), ".shareit")
}

func (*FileVault) New() *FileVault {
	return &FileVault{
		fileMap: make(map[string]filePathDetails),
	}
}

var fv *FileVault

func main() {
	fv = fv.New()
	go establishPipe()
	StartServer(DEFAULT_SERVER_PORT, SERVER_FILE)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	log.Printf("[SERVER]: Recieved token: %s\n", tokenString)

	// Parsing unverified token to get payload
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		log.Printf("[SERVER]: %s\n", err)
		return
	}

	// Extracting filepath from token payload
	var filepath string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		filepath = fmt.Sprintf("%s", claims["filepath"])
	} else {
		http.Error(w, "Invalid token!", http.StatusUnauthorized)
		return
	}

	log.Printf("[SERVER]: filepath corresponding to given token: %s\n", filepath)

	// Extracting JWT token secret for the corrsponding filepath from FileVault
	fileDetails, ok := fv.Read(filepath)
	if !ok {
		http.Error(w, "Invalid token payload!", http.StatusUnauthorized)
		return
	}
	secret, _ := base32.StdEncoding.DecodeString(fileDetails.Secret)

	// Verifying token with secret
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		http.Error(w, "Invalid token, unable to parse!", http.StatusUnauthorized)
		log.Printf("[SERVER]: %s\n", err)
		return
	}

	if !token.Valid {
		http.Error(w, "Unauthorised token!", http.StatusUnauthorized)
		return
	}

	rangeHeader := r.Header.Get("Range")

	var start, end int64

	if filepath == "" {
		http.Error(w, "filepath param misssing!", http.StatusBadRequest)
		return
	}

	// Checking if the file exists
	// If yes then proceed with downloading
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

// Unix socket connection handler
// Recieves a socket payload and writes it to FileVault
func handleSocketConnection(conn net.Conn) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)

	var payload socketPayload
	err := dec.Decode(&payload)

	if err != nil {
		log.Printf("[SERVER SOCKET]: Unable to decode socket: %s\n", err)
	}

	fv.Write(payload.FilePath, payload.FileAuth)

	log.Printf("[SERVER SOCKET]: Payload received!")

}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pong!")
}

func StartServer(port string, serverFilePath string) {
	log.Print("Starting server...")

	err := os.WriteFile(serverFilePath, []byte(port), 0644)

	if err != nil {
		log.Fatalf("[SERVER]: Unable to write server details: %s\n", err)
	}

	log.Printf("Server started on port: %s\n", port)

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/", handleFile)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func isValidPath(filePath string) bool {
	_, err := os.Stat(filePath)

	if err != nil {
		return false
	}

	return true
}

func establishPipe() {
	// Establish unix domain socket
	sockAddr, err := filepath.Abs(UNIX_SOCKET_FILE)

	if err != nil {
		log.Fatalf("[SERVER]: Unable to resolve socket file address: %s\n", err)
	}

	// Remove any existing socket file with same name
	_ = os.Remove(sockAddr)

	log.Printf("[SERVER]: Connecting to socket @: %s\n", sockAddr)

	listener, err := net.Listen("unix", sockAddr)

	if err != nil {
		log.Fatalf("[SERVER]: Error creating socket listener")
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("[SERVER]: Error accepting unix socket connection: %s\n", err)
		}

		go handleSocketConnection(conn)
	}
}
