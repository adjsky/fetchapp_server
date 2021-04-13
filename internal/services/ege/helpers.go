package ege

import (
	"os/exec"
	"strconv"
	"strings"
)

const pythonScriptPath string = "internal/services/ege/python/main.py"

func executeScript(args ...string) (string, error) {
	var out strings.Builder
	var errOut strings.Builder
	command := exec.Command("python", args...)
	command.Stdout = &out
	command.Stderr = &errOut
	err := command.Run()
	if err != nil {
		return errOut.String(), err
	}
	return out.String(), nil
}

func processQuestion(questionNumber int, filepath string, req *question24Request) (string, error) {
	return executeScript(pythonScriptPath, "solve", strconv.Itoa(questionNumber), "-f",
		filepath, "-t", strconv.Itoa(req.Type), "-c", req.Char)
}
