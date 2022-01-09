package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Success case simple", func(t *testing.T) {
		cmd := []string{"pwd", "-L"}
		exitCode := RunCmd(cmd, Environment{})

		require.Equal(t, exitCodeOk, exitCode)
	})

	t.Run("Success case with filled dir env", func(t *testing.T) {
		cmd := []string{"pwd", "-L"}
		exitCode := RunCmd(cmd, Environment{
			"TEST_QWE": EnvValue{
				Value: "TEST_QWE",
			},
		})

		require.Equal(t, exitCodeOk, exitCode)
		require.Contains(t, os.Environ(), "TEST_QWE=TEST_QWE")
	})

	t.Run("Exec err case", func(t *testing.T) {
		cmd := []string{"pwd", "-KEK"}
		exitCode := RunCmd(cmd, Environment{
			"TEST_QWE": EnvValue{
				Value: "TEST_QWE",
			},
		})

		require.Equal(t, 1, exitCode)
	})
}

func TestFillEnv(t *testing.T) {
	t.Run("Success simple case", func(t *testing.T) {
		dirEnv := Environment{
			"A": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
			"TEST_B": EnvValue{
				Value:      "B",
				NeedRemove: false,
			},
			"TEST_C": EnvValue{
				Value:      "123",
				NeedRemove: false,
			},
		}
		expectedRes := []string{"TEST_B=B", "TEST_C=123"}

		resCode := fillEnv(dirEnv)

		require.Equal(t, 0, resCode)
		require.Contains(t, os.Environ(), expectedRes[0])
		require.Contains(t, os.Environ(), expectedRes[1])
	})
}
