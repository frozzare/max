package exec

import (
	"bufio"
	"os"
	goexec "os/exec"
	"strings"
)

// Cmd will execute a input cmd string.
func Cmd(input string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	parts := strings.Fields(input)
	head := parts[0]
	parts = parts[1:]

	cmd := goexec.Command(head, parts...)
	cmd.Dir = path

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	var res []string

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			res = append(res, scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return strings.Join(res, "\n"), nil
}
