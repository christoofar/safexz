// Description: This file contains the implementation of the decompression functions for safexz.
// These functions are used to decompress data that has been compressed using the xz format.
// You do not need to worry about the C language bindings, as they are handled in the lzma package.
package safexz

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	internal "github.com/christoofar/safexz/internal"
)

func DecompressString(compressedString string) (string, error) {
	return "", nil
}

func DecompressBytes(compressedBytes []byte) ([]byte, error) {

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	internal.DecompressIn(readchan, writechan)

	go func() {
		for i := 0; i < len(compressedBytes); i += internal.MAX_BUF_SIZE {
			end := i + internal.MAX_BUF_SIZE
			if end > len(compressedBytes) {
				end = len(compressedBytes)
			}
			readchan <- compressedBytes[i:end]
		}
		close(readchan)
	}()

	var decompressed []byte
	for data := range writechan {
		decompressed = append(decompressed, data...)
	}
	return decompressed, nil

}

func DecompressFile(input_path, output_path string) error {
	return nil
}

func DecompressFileWithProgress(inpath, outpath string, progress func(uint64, uint64)) error {

	// Check the file extension
	extension := filepath.Ext(inpath)
	fileExtension := extension[1:]
	if fileExtension != "xz" {
		return fmt.Errorf("the input file [%s] should probably have an xz extension, can you go look?", outpath)
	}

	f, err := os.Open(inpath)
	if err != nil {
		return err
	}

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	internal.DecompressIn(readchan, writechan)
	var readCount uint64
	var writeCount uint64

	readfunc := func() {
		readbuf := make([]byte, internal.MAX_BUF_SIZE)

		for {
			bytes, err := f.Read(readbuf)
			readCount += uint64(bytes)
			if progress != nil && readCount%16550 == 0 {
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

	// If the outpath exists, delete it
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

// DecompressFileToMemory reads a file from the filesystem and decompresses it into memory as a byte slice buffer.
// This is useful for small files that you want to decompress and then work with in memory, such as scanning compressed logs.
func DecompressFileToMemory(path string) ([]byte, error) {

		// Check the file extension
		extension := filepath.Ext(path)
		fileExtension := extension[1:]
		if fileExtension != "xz" {
			return []byte{}, fmt.Errorf("the input file [%s] should probably have an xz extension, can you go look?", path)
		}

		f, err := os.Open(path)
		if err != nil {
			return []byte{}, err
		}

	inputchan := make(chan []byte, 1)
	outputchan := make(chan []byte, 1)

	go func() {
		internal.DecompressIn(inputchan, outputchan)
	}()

	readbuf := make([]byte, internal.MAX_BUF_SIZE)
	for {
		bytes, err := f.Read(readbuf)
		if err != nil {
			close(inputchan)
			break
		}
		inputchan <- readbuf[:bytes]
	}

	outputbuf := bytes.Buffer{}
	for data := range outputchan {
		_, err := outputbuf.Write(data)
		if err != nil {
			return []byte{}, err
		}
	}

	return outputbuf.Bytes(), nil
}

// DecompressStream skips a call to io.Copy() by just compressing whatever stream you put in the
// input reader and writing it to the output writer.
func DecompressStream(input io.Reader, output io.Writer) error {
	inputchan := make(chan []byte, 1)
	outputchan := make(chan []byte, 1)

	go func() {
		internal.DecompressIn(inputchan, outputchan)
	}()

	readbuf := make([]byte, internal.MAX_BUF_SIZE)
	for {
		bytes, err := input.Read(readbuf)
		if err != nil {
			close(inputchan)
			break
		}
		inputchan <- readbuf[:bytes]
	}

	for data := range outputchan {
		_, err := output.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
