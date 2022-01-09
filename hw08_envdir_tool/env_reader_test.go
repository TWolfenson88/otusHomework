package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	emptyDirPath = "testdata/empty_dir"

	filledParentDirPath = "testdata/parent_dir"
	childDirPath        = "testdata/parent_dir/child_dir"

	fileWithBadNameDirPath = "testdata/bad_name_test"
	fileWithBadNamePath    = "testdata/bad_name_test/BAD=NAME"

	successCaseDirPath = "testdata/env"
)

func TestReadDir(t *testing.T) {
	defer func() {
		os.Remove(emptyDirPath)

		os.RemoveAll(filledParentDirPath)

		os.RemoveAll(fileWithBadNameDirPath)
	}()

	t.Run("Empty dir", func(t *testing.T) {
		err := os.Mkdir(emptyDirPath, 0o755)
		if err != nil {
			log.Fatal(err)
		}

		res, err := ReadDir(emptyDirPath)
		require.NoError(t, err)
		require.Len(t, res, 0)
	})

	t.Run("Dir with only other dir inside", func(t *testing.T) {
		err := os.Mkdir(filledParentDirPath, 0o755)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Mkdir(childDirPath, 0o755)
		if err != nil {
			log.Fatal(err)
		}

		res, err := ReadDir(filledParentDirPath)

		require.Error(t, err)
		require.EqualError(t, err, ErrFileIsADirectory.Error())
		require.Len(t, res, 0)
	})

	t.Run("Bad file name err", func(t *testing.T) {
		err := os.Mkdir(fileWithBadNameDirPath, 0o755)
		if err != nil {
			log.Fatal(err)
		}

		_, err = os.Create(fileWithBadNamePath)
		if err != nil {
			log.Fatal(err)
		}

		res, err := ReadDir(fileWithBadNameDirPath)

		require.Error(t, err)
		require.EqualError(t, err, ErrBadFileName.Error())
		require.Len(t, res, 0)
	})

	t.Run("Success case", func(t *testing.T) {
		expectedRes := Environment{
			"BAR": EnvValue{
				Value: "bar",
			},
			"EMPTY": EnvValue{
				NeedRemove: true,
			},
			"FOO": EnvValue{
				Value: "   foo\nwith new line",
			},
			"HELLO": EnvValue{
				Value: `"hello"`,
			},
			"UNSET": EnvValue{
				NeedRemove: true,
			},
		}

		res, err := ReadDir(successCaseDirPath)

		require.NoError(t, err)
		require.Equal(t, expectedRes, res)
	})
}
