// Description: This file contains the implementation of the compression functions.  These
// functions are used to compress data using the xz format.  You don't need to worry about
// the C language bindings, as they are handled in the lzma package.
package safexz

import (
	"io"
	"os"

	internal "github.com/christoofar/safexz/internal/common"
)

func CompressString(s string) (string, error) {
	return "", nil
}

func CompressBytes(b []byte) ([]byte, error) {
	return nil, nil
}

func CompressFile(path string) error {
	return nil
}

func CompressFileWithProgress(path string, progress func(uint64)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	readchan := make(chan []byte)
	writechan := make(chan []byte)

	readbuf := make([]byte, internal.MAX_BUF_SIZE)
	internal.CompressIn(&readchan, &writechan)
	go func() {
		for {
			bytes, err := f.Read(readbuf)
			if err != nil {
				close(readchan)
				break
			}
			readchan <- readbuf[:bytes]
		}
	}()

	outfile, _ := os.Create("output.xz")
	var count uint64
	for data := range writechan {
		outfile.Write(data)
		count += uint64(len(data))
		if progress != nil {
			progress(count)
		}
	}
	outfile.Close()

	return nil
}

func CompressFileToMemory(path string) ([]byte, error) {
	return nil, nil
}

func CompressStream(input io.Reader, output io.Writer) error {
	return nil
}
