package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/skip2/go-qrcode"
)

type filePathDetails struct {
	Secret string
	Otp    string
}

type socketPayload struct {
	FileAuth filePathDetails
	FilePath string
}

const DEFAULT_SERVER_PORT = "8966"
const SECRET_LENGTH = 32

var RUNTIME_DIR = getRuntimeDirectory()
var SERVER_FILE = filepath.Join(RUNTIME_DIR, "server.pid")
var UNIX_SOCKET_FILE = filepath.Join(RUNTIME_DIR, "server.sock")
var SERVER_LOG_FILE_PATH = filepath.Join(RUNTIME_DIR, "server.log")

func getRuntimeDirectory() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "ShareIT")
	}

	return filepath.Join(os.Getenv("HOME"), ".shareit")
}

func main() {
	filepath := flag.String("filepath", "", "Absolute address of file")
	otp := flag.String("otp", "", "One time password for file sharing")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *filepath == "" {
		fmt.Fprintf(os.Stderr, "Error: -filepath flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if isValidPath(*filepath) {
		// Get Local IP from network interfaces
		ips := getLocalIP()

		if len(ips) == 0 {
			log.Fatalf("No LAN detected!")
		}

		// Get port for the server from server.pid
		port, err := getServerPort()
		retriedConnection := false

		if err != nil {
			log.Println("Server not running!")

			// Start server process in case no saved port is found
			startServerProcess()
			port = DEFAULT_SERVER_PORT
			retriedConnection = true
		} else {
			log.Printf("Saved port: %s\n", port)

			// Try starting server process in case tcp dialup to server fails
			// on saved port
			if !isServerUp(port) {
				log.Println("Server connection failed, retrying...")
				startServerProcess()
				retriedConnection = true
			} else {
				log.Printf("Server connection successful on port: %s\n", port)
			}
		}

		if retriedConnection {
			if !isServerUp(port) {
				log.Fatal("Server connection failed!")
			} else {
				log.Printf("Server connection successful on port: %s\n", port)
			}
		}

		// Generate random token secret, secret is base32 encoded string
		// which will also be shared with server porcess via unix socket
		secretString, err := generateFilePathSecret(SECRET_LENGTH)
		if err != nil {
			log.Fatalln("[CLI]: Unable to generate token secret!!")
		}

		// Generate token with filepath as payload
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"filepath": *filepath,
		})

		// Converting secret to original bytes and signing token with it
		secret, _ := base32.StdEncoding.DecodeString(secretString)
		token, err := t.SignedString(secret)
		if err != nil {
			log.Fatalf("[CLI]: Unable to generate token: %s\n", err)
		}

		// Creating unix socket payload in common acceptable format i.e. socketPayload
		payload := socketPayload{
			FileAuth: filePathDetails{
				Otp:    *otp,
				Secret: secretString,
			},
			FilePath: *filepath,
		}

		// Send secret and password to server, along with filepath
		sendFilePayload(payload)

		// LINK generation
		params := url.Values{}
		params.Set("token", token)
		encodedParams := params.Encode()

		for _, ip := range ips {
			generateQRCode(ip.String(), port, "?"+encodedParams)
		}
	} else {
		log.Fatalf("Invlaid file path: %s\n", *filepath)
	}
}

// TCP dialup to server process to check server's running status
func isServerUp(port string) bool {
	_, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Printf("Server not up: %s\n", err)
		return false
	}

	return true
}
func generateQRCode(ip string, port string, path string) {
	link := fmt.Sprintf("http://%s:%s/%s", ip, port, path)

	qrCode, err := qrcode.New(link, qrcode.Medium)

	if err != nil {
		log.Fatalf("Error in generating qrcode: %s\n", err)
	}

	fmt.Printf("\n%s\n", qrCode.ToSmallString(false))

	fmt.Printf("Link: %s\n", link)
}

func getServerPort() (string, error) {
	data, err := os.ReadFile(SERVER_FILE)
	if err != nil {
		return "", err
	}

	if _, err = strconv.ParseUint(string(data), 10, 64); err != nil {
		return "", err
	}

	return string(data), nil
}

func isValidPath(filePath string) bool {
	if !filepath.IsAbs(filePath) {
		log.Println("Absolute path required!")
		return false
	}

	stat, err := os.Stat(filePath)

	if err != nil {
		return false
	}

	log.Println("File Stats: ")
	log.Printf("File Name: %s\n", stat.Name())
	log.Printf("File Size: %d KB\n\n", stat.Size()/1024)

	return true
}

func getLocalIP() []net.IP {
	interfaces, err := net.Interfaces()
	var ips []net.IP

	if err != nil {
		log.Fatalln("No network interfaces found!")
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagBroadcast != 0 && i.Flags&net.FlagLoopback == 0 && i.Flags&net.FlagRunning != 0 {
			addresses, err := i.Addrs()

			if err != nil {
				continue
			}

			var localIP net.IP

			for _, address := range addresses {
				switch v := address.(type) {
				case *net.IPAddr:
					if v.IP.To4() != nil {
						localIP = v.IP.To4()
					}
				case *net.IPNet:
					if v.IP.To4() != nil {
						localIP = v.IP.To4()
					}
				}

				if localIP != nil {
					ips = append(ips, localIP)
				}
			}
		}
	}
	return ips
}

/*
Runs server executable file as a separate process
*/
func startServerProcess() {
	serverExecPath := filepath.Join(".", fmt.Sprintf("shareit.server.%s", runtime.GOOS))

	if runtime.GOOS == "windows" {
		serverExecPath += ".exe"
	}

	log.Printf("Starting server process: %s\n", serverExecPath)

	// Create/Open server log file in append mode from SERVER_LOG_FILE_PATH
	file, err := os.OpenFile(SERVER_LOG_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatalf("Unable to create/open server log file: %s\n", err)
	}

	cmd := exec.Command(serverExecPath)
	cmd.Stdin = file
	cmd.Stdout = file
	cmd.Stderr = file

	err = cmd.Start()

	if err != nil {
		log.Fatalf("Unable to start server process: %s\n", err)
	}
	log.Println("Server started")
	time.Sleep(1 * time.Second)
}

func generateFilePathSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(bytes), nil
}

func sendFilePayload(payload socketPayload) {
	absSocketAddres, err := filepath.Abs(UNIX_SOCKET_FILE)

	if err != nil {
		log.Fatalf("[CLI]: Unable to resolve socket file address: %s\n", err)
	}

	addr, err := net.ResolveUnixAddr("unix", absSocketAddres)

	if err != nil {
		log.Fatalf("[CLI]: Unable to resolve unix address: %s\n", err)
	}

	conn, err := net.DialUnix("unix", nil, addr)

	if err != nil {
		log.Fatalf("[CLI]: Unable to connect to socket: %s\n", err)
	}

	// Encoding struct payload to socket
	enc := gob.NewEncoder(conn)

	err = enc.Encode(payload)

	if err != nil {
		log.Fatalf("[CLI]: Unable to encode payload to socket: %s\n", err)
	}
}
