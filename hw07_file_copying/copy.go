package main

import (
	"errors"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const MB = 1024 * 1024

func Copy(fromPath, toPath string, offset, limit int64) error {
	fStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if !fStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fSize := fStat.Size()
	if fSize <= offset {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = fSize - offset
	}

	in, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = in.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bytesCopy := limit
	readSize := int64(math.Min(float64(MB), float64(limit)))
	var totalReadBytes int64

	bar := pb.Full.Start64(bytesCopy)
	fToBarProxy := bar.NewProxyWriter(out)

	for {
		n, err := io.CopyN(fToBarProxy, in, readSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		totalReadBytes += n
		if totalReadBytes >= bytesCopy {
			break
		}
	}

	bar.Finish()

	return nil
}
