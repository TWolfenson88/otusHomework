package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrBadFileName      = errors.New("file name contains '=' symbol")
	ErrFileIsADirectory = errors.New("file is a directory")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envVars := make(Environment, len(filesInfo))
	for _, fi := range filesInfo {
		if fi.IsDir() {
			return nil, ErrFileIsADirectory
		}

		if strings.Contains(fi.Name(), "=") {
			return nil, ErrBadFileName
		}

		file, err := os.Open(filepath.Join(dir, fi.Name()))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		envVars[fi.Name()] = filterEnvVal(scanner.Text())
	}

	return envVars, nil
}

func filterEnvVal(s string) EnvValue {
	var envVal EnvValue

	s = strings.TrimRightFunc(s, unicode.IsSpace)
	s = strings.ReplaceAll(s, string([]byte{'\x00'}), "\n")
	if len(s) == 0 {
		envVal.NeedRemove = true
	}

	envVal.Value = s

	return envVal
}