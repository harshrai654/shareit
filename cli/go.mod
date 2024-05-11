module ttharsh.shareit/cli

go 1.22.1

replace ttharsh.shareit/server => ../server

require ttharsh.shareit/server v0.0.0-00010101000000-000000000000

require (
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/nsf/termbox-go v1.1.1 // indirect
	github.com/yeqown/go-qrcode/v2 v2.2.4 // indirect
	github.com/yeqown/go-qrcode/writer/terminal v1.1.1 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
)
