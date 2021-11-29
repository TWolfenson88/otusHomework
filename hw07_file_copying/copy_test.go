package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	inputFile  = "./testdata/input.txt"
	outputFile = "./testdata/output.txt"
)

func readFromFiles(t *testing.T, expFile string) ([]byte, []byte) {
	t.Helper()
	resultText, err := ioutil.ReadFile(outputFile)
	require.NoError(t, err)

	expectedText, err := ioutil.ReadFile(expFile)
	require.NoError(t, err)

	return resultText, expectedText
}

func TestCopy(t *testing.T) {
	t.Cleanup(func() {
		err := os.Remove(outputFile)
		require.NoError(t, err)
	})

	t.Run("offset 0, limit 0", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 0, 0)
		require.NoError(t, err)

		resultText, expectedText := readFromFiles(t, "./testdata/out_offset0_limit0.txt")

		require.Equal(t, expectedText, resultText)
	})

	t.Run("offset 0, limit 1000", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 0, 1000)
		require.NoError(t, err)

		resultText, expectedText := readFromFiles(t, "./testdata/out_offset0_limit1000.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("offset 100, limit 1000", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 100, 1000)
		require.NoError(t, err)

		resultText, expectedText := readFromFiles(t, "./testdata/out_offset100_limit1000.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("offset 6000, limit 1000", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 6000, 1000)
		require.NoError(t, err)

		resultText, expectedText := readFromFiles(t, "./testdata/out_offset6000_limit1000.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("Error: offset more than file size", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 0, 0)
		require.NoError(t, err)

		require.Error(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("Error: file with zero size", func(t *testing.T) {
		err := Copy("/dev/urandom", outputFile, 0, 0)

		require.Error(t, ErrUnsupportedFile, err)
	})

	t.Run("Limit more than file. Read until EOF", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 0, 9999999999)
		require.NoError(t, err)

		resultText, expectedText := readFromFiles(t, "./testdata/out_offset0_limit0.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("Should return nil if all are ok", func(t *testing.T) {
		err := Copy(inputFile, outputFile, 0, 0)

		require.Nil(t, err)
	})
}
