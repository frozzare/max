package exec

import (
	"bufio"
	"os"
	goexec "os/exec"
	"regexp"
	"strings"
)

func setVariables(input string) (string, error) {
	re := regexp.MustCompile(`(\w+\=\w+)`)
	match := re.FindAllStringSubmatch(input, -1)

	for _, row := range match {
		if len(row) < 1 {
			continue
		}

		input = strings.Replace(input, row[1], "", -1)
		p := strings.Split(row[1], "=")

		if err := os.Setenv(p[0], p[1]); err != nil {
			return input, err
		}
	}

	return strings.TrimLeft(input, " "), nil
}

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

	input, err := setVariables(input)
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
		return strings.Join(res, "\n"), err
	}

	if err := cmd.Wait(); err != nil {
		return strings.Join(res, "\n"), err
	}

	return strings.Join(res, "\n"), nil
}
