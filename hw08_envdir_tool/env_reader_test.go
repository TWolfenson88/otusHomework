package main

import (
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

//
//func TestReadDir(t *testing.T) {
//	dirName := "./testdata/env"
//	expectedEnvs := Environment{
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
//	t.Run("positive tests", func(t *testing.T) {
//		t.Run("read envs from dir", func(t *testing.T) {
//			envs, err := ReadDir(dirName)
//
//			require.Equal(t, expectedEnvs, envs)
//			require.NoError(t, err)
//		})
//
//		t.Run("use only files names and ignore dirs", func(t *testing.T) {
//			dirIgnoreName := "./testdata/env/ingnoredir"
//			os.Mkdir(dirIgnoreName, os.ModeDir)
//			defer os.RemoveAll(dirIgnoreName)
//
//			envs, err := ReadDir(dirName)
//
//			require.NotContains(t, envs, "ingnoredir")
//			require.NoError(t, err)
//		})
//	})
//
//	t.Run("negative tests", func(t *testing.T) {
//		t.Run("read envs from not existing dir", func(t *testing.T) {
//			envs, err := ReadDir("./notExistingDir")
//
//			require.Nil(t, envs)
//			require.IsType(t, new(os.PathError), errors.Cause(err))
//		})
//
//		t.Run("read envs from from dir with '=' in the name", func(t *testing.T) {
//			fileIncorrectName := "./testdata/env/FT=42"
//			os.Create(fileIncorrectName)
//			defer os.Remove(fileIncorrectName)
//
//			envs, err := ReadDir("./testdata/env")
//
//			require.Nil(t, envs)
//			require.IsType(t, ErrIncorrectFileName, errors.Cause(err))
//		})
//	})
//}


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
		err := os.Mkdir(emptyDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		res, err := ReadDir(emptyDirPath)
		require.NoError(t, err)
		require.Len(t, res, 0)
	})

	t.Run("Dir with only other dir inside", func(t *testing.T) {
		err := os.Mkdir(filledParentDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Mkdir(childDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		res, err := ReadDir(filledParentDirPath)

		require.Error(t, err)
		require.EqualError(t, err, ErrFileIsADirectory.Error())
		require.Len(t, res, 0)
	})

	t.Run("Bad file name err", func(t *testing.T) {
		err := os.Mkdir(fileWithBadNameDirPath, 0755)
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
