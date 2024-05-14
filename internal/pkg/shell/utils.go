package shell

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func SplitArgs(args string) []string {
	return strings.Split(args, " ")
}

func MkDir(path, dir string) {
	os.Mkdir(path+dir, os.ModePerm)
}

func Exec(command, args string) (string, string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, SplitArgs(args)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func CountFiles(dir string) string {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("ls", "-1", dir, "|", "wc", "-l")

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// err := cmd.Run()
	cmd.Run()
	return stdout.String()
}
