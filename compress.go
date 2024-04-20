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
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	internal.CompressIn(readchan, writechan, int(use_strategy))
	go func() {
		for i := 0; i < len(s); i += internal.MAX_BUF_SIZE {
			end := i + internal.MAX_BUF_SIZE
			if end > len(s) {
				end = len(s)
			}
			readchan <- []byte(s[i:end])
		}
		close(readchan)
	}()

	var compressed string
	for data := range writechan {
		compressed += string(data)
	}

	return compressed, nil
}

func CompressBytes(b []byte, strategy ...CompressionStrategy) ([]byte, error) {
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	internal.CompressIn(readchan, writechan, int(use_strategy))

	go func() {
		for i := 0; i < len(b); i += internal.MAX_BUF_SIZE {
			end := i + internal.MAX_BUF_SIZE
			if end > len(b) {
				end = len(b)
			}
			readchan <- b[i:end]
		}
		close(readchan)
	}()

	var compressed []byte
	for data := range writechan {
		compressed = append(compressed, data...)
	}

	return compressed, nil
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

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	internal.CompressIn(readchan, writechan, int(use_strategy))
	var readCount uint64
	var writeCount uint64

	readfunc := func() {
		readbuf := make([]byte, internal.MAX_BUF_SIZE)

		for {
			bytes, err := f.Read(readbuf)
			readCount += uint64(bytes)
			if progress != nil && readCount%4096 == 0 {
				progress(readCount, writeCount)
			}
			if err != nil { // The EOF has been hit, send the final batch
				readchan <- readbuf[:bytes]
				close(readchan)
				break
			}

			data := make([]byte, bytes)
			copy(data, readbuf)
			readchan <- data
		}
	}

	// If the outfpath exists, delete it
	if _, err := os.Stat(outpath); err == nil {
		err := os.Remove(outpath)
		if err != nil {
			return err
		}
	}
	outfile, err := os.Create(outpath)
	if err != nil {
		return err
	}

	go readfunc()

	donewrite := make(chan bool, 1)
	go func() {
		for data := range writechan {
			outfile.Write(data)
			if len(data) > 0 {
				if progress != nil {
					writeCount += uint64(len(data))
					progress(readCount, writeCount)
				}
			}
		}
		donewrite <- true
	}()
	<-donewrite
	outfile.Close()

	return nil
}

func CompressFileToMemory(path string) ([]byte, error) {
	return nil, nil
}

func CompressStream(input io.Reader, output io.Writer) error {
	return nil
}
