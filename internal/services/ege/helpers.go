package ege

import (
	"os/exec"
	"strconv"
	"strings"
)

func executeScript(args ...string) (string, error) {
	var out strings.Builder
	var errOut strings.Builder
	command := exec.Command("python", args...)
	command.Stdout = &out
	command.Stderr = &errOut
	if err := command.Run(); err != nil {
		return errOut.String(), err
	}
	return out.String(), nil
}

func processQuestion(scriptPath string, questionNumber int, filepath string, req *question24Request) (string, error) {
	return executeScript(scriptPath, "solve", strconv.Itoa(questionNumber), "-f",
		filepath, "-t", strconv.Itoa(req.Type), "-c", req.Char)
}
