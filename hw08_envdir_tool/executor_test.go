package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

//
//func TestRunCmd(t *testing.T) {
//	envs := Environment{
//		"BAR": EnvValue{
//			Value:      "bar",
//			NeedRemove: false,
//		},
//		"EMPTY": EnvValue{
//			Value:      "",
//			NeedRemove: false,
//		},
//		"FOO": EnvValue{
//			Value:      "   foo\nwith new line",
//			NeedRemove: false,
//		},
//		"HELLO": EnvValue{
//			Value:      "\"hello\"",
//			NeedRemove: false,
//		},
//		"UNSET": EnvValue{
//			Value:      "",
//			NeedRemove: true,
//		},
//	}
//
//	t.Run("test with argument", func(t *testing.T) {
//		cmd := []string{"./testdata/echo.sh", "foo=42"}
//		exitCode := RunCmd(cmd, envs)
//
//		require.Zero(t, exitCode)
//	})
//
//	t.Run("test without arguments", func(t *testing.T) {
//		cmd := []string{"./testdata/echo.sh"}
//		exitCode := RunCmd(cmd, envs)
//
//		require.Zero(t, exitCode)
//	})
//
//	t.Run("test without arguments", func(t *testing.T) {
//		cmd := []string{"./test.sh"}
//		exitCode := RunCmd(cmd, envs)
//
//		require.Zero(t, exitCode)
//	})
//
//	t.Run("handle correct exit code", func(t *testing.T) {
//		cmd := []string{"./testdata/exit42.sh"}
//		exitCode := RunCmd(cmd, envs)
//
//		require.Equal(t, 42, exitCode)
//	})
//}


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
