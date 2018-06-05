package docker

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func generateScript(commands []string) string {
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(commands, "\n")))
}

func generateCommand(command string) string {
	return fmt.Sprintf("echo %s | base64 -d | /bin/sh -e", command)
}
