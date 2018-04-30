package main

import "github.com/frozzare/max/cmd"

var version = "master"

func main() {
	cmd.Version = version
	cmd.Execute()
}
