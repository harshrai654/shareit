package main

import (
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

	"github.com/skip2/go-qrcode"
)

const SERVER_FILE = "./server.pid"
const DEFAULT_SERVER_PORT = "8965"

func main() {

	filepath := flag.String("filepath", "", "Absolute address of file")
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
		// Get Local IP
		ips := getLocalIP()

		if len(ips) == 0 {
			log.Fatalf("No LAN detected!")
		}

		// Get port for the server
		port, err := getServerPort()
		retriedConnection := false

		if err != nil {
			log.Println("Server not running!")
			// Try starting server
			startServerProcess()
			port = DEFAULT_SERVER_PORT
			retriedConnection = true
		} else {
			log.Printf("Saved port: %s\n", port)
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

		params := url.Values{}
		params.Set("path", *filepath)
		encodedParams := params.Encode()

		for _, ip := range ips {
			generateQRCode(ip.String(), port, "?"+encodedParams)
		}
	} else {
		log.Fatalf("Invlaid file path: %s\n", *filepath)
	}
}

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

func startServerProcess() {
	serverExecPath := filepath.Join(".", fmt.Sprintf("shareit.server.%s", runtime.GOOS))

	if runtime.GOOS == "windows" {
		serverExecPath += ".exe"
	}

	cmd := exec.Command(serverExecPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		log.Fatalf("Unable to start server process: %s\n", err)
	}
	log.Println("Server started")
	time.Sleep(1 * time.Second)
}
