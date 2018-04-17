package exec

import (
	"bufio"
	"os"
	goexec "os/exec"
	"strings"
)

// Cmd will execute a input cmd string.
func Cmd(input string, args ...string) (string, error) {
	var path string

	if len(args) > 0 {
		path = args[0]
	}

	if len(path) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		path = wd
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
		return strings.Join(res, "\n"), err
	}

	if err := cmd.Wait(); err != nil {
		return strings.Join(res, "\n"), err
	}

	return strings.Join(res, "\n"), nil
}
