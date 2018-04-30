package cmd

import (
	"log"

	"github.com/frozzare/max/internal/config"
	"github.com/spf13/pflag"
)

var (
	// Version of max
	Version = "master"
)

func runCommands(cmd string, args []string) bool {
	switch cmd {
	case "cache":
		if len(args) == 0 || args[0] != "flush" {
			return false
		}

		if cache, err := config.CreateCache(); err == nil {
			cache.Flush()
		}

		log.Println("max: cache flushed")

		return true
	case "help":
		if len(args) > 0 {
			return false
		}

		pflag.Usage()
		return true
	case "version":
		log.Printf("max version %s\n", Version)
		return true
	default:
		return false
	}
}
