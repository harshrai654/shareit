package main

import "ttharsh.shareit/server/lib"

const DEFAULT_SERVER_PORT = "8965"
const SERVER_FILE = "../server.pid"

func main() {
	lib.StartServer(DEFAULT_SERVER_PORT, SERVER_FILE)
}
