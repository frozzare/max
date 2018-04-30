package main

import (
	"github.com/frozzare/max/cmd"
	_ "github.com/joho/godotenv/autoload"
)

var version = "master"

func main() {
	cmd.Version = version
	cmd.Execute()
}
