package main

import (
	"embed"
	"github.com/uouuou/ServerManagerSystem/cmd"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
)

//go:embed web
var f embed.FS

func main() {
	mid.FS = f
	cmd.Install()
}
