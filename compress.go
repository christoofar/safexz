// Description: This file contains the implementation of the compression functions.  These
// functions are used to compress data using the xz format.  You don't need to worry about
// the C language bindings, as they are handled in the lzma package.
package safexz

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	internal "github.com/christoofar/safexz/internal/common"
)

func CompressString(s string, strategy ...CompressionStrategy) (string, error) {
	return "", nil
}

func CompressBytes(b []byte, strategy ...CompressionStrategy) ([]byte, error) {
	return nil, nil
}

func CompressFile(inpath string, outpath string, strategy ...CompressionStrategy) error {
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}
	return CompressFileWithProgress(inpath, outpath, nil, use_strategy)
}

func CompressFileWithProgress(inpath string, outpath string, progress func(uint64, uint64), strategy ...CompressionStrategy) error {

	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}

	// Check the file extension
	extension := filepath.Ext(outpath)
	fileExtension := extension[1:]
	if fileExtension != "xz" {
		return fmt.Errorf("output file [%s] must have an xz extension", outpath)
	}

	f, err := os.Open(inpath)
	if err != nil {
		return err
	}

	readchan := make(chan []byte)
	writechan := make(chan []byte)

	readbuf := make([]byte, internal.MAX_BUF_SIZE)
	internal.CompressIn(&readchan, &writechan, use_strategy)
	var readCount uint64
	var writeCount uint64

	go func() {
		for {
			bytes, err := f.Read(readbuf)
			readCount += uint64(bytes)
			if progress != nil && readCount%4096 == 0 {
				progress(readCount, writeCount)
			}
			if err != nil {
				close(readchan)
				break
			}
			readchan <- readbuf[:bytes]
		}
	}()

	outfile, err := os.Create(outpath)
	if err != nil {
		return err
	}

	for data := range writechan {
		outfile.Write(data)
		if len(data) > 0 {
			if progress != nil {
				writeCount += uint64(len(data))
				progress(readCount, writeCount)
			}
		}
	}
	outfile.Close()
	readbuf = nil

	return nil
}

func CompressFileToMemory(path string) ([]byte, error) {
	return nil, nil
}

func CompressStream(input io.Reader, output io.Writer) error {
	return nil
}
