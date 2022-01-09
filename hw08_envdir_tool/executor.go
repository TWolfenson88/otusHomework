package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	exitCodeOk             = 0
	exitCodeUnknown        = 400
	exitCodeCannotUnsetEnv = 401
	exitCodeCannotSetEnv   = 402
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	exCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	code := fillEnv(env)
	if code != 0 {
		return code
	}

	exCmd.Stdout = os.Stdout
	exCmd.Stderr = os.Stderr
	if err := exCmd.Run(); err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return exitCodeUnknown
	}

	return exitCodeOk
}

func fillEnv(env Environment) int {
	for key, elem := range env {
		if elem.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return exitCodeCannotUnsetEnv
			}

			continue
		}

		err := os.Setenv(key, elem.Value)
		if err != nil {
			return exitCodeCannotSetEnv
		}
	}

	return exitCodeOk
}
