package lib

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

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

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
