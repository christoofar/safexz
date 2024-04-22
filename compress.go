// Description: This file contains the implementation of the compression functions.  These
// functions are used to compress data using the xz format.  You don't need to worry about
// the C language bindings, as they are handled in the lzma package.
package safexz

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	/* trunk-ignore(golangci-lint/typecheck) */
	internal "github.com/christoofar/safexz/internal"
)

// CompressString compresses a string using the xz format and returns the compressed string.
func CompressString(s string, strategy ...CompressionStrategy) (string, error) {
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	var funcerr error = nil
	go func() {
		funcerr = internal.CompressIn(readchan, writechan, int(use_strategy))
	}()
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

	return compressed, funcerr
}

// CompressBytes compresses a byte slice using the xz format and returns the compressed byte slice.  If the byte slice is huge,
// you may want to consider using CompressFile or CompressStream instead.  The reason is that the compression process can greatly expand the
// amount of memory consumed depending on the CompressionStrategy used.
// The compression process can greatly expand the amount of memory consumed depending on the CompressionStrategy used.
func CompressBytes(b []byte, strategy ...CompressionStrategy) ([]byte, error) {
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	var funcerr error = nil
	go func() {
		funcerr = internal.CompressIn(readchan, writechan, int(use_strategy))
	}()

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

	return compressed, funcerr
}

// CompressFile compresses a file using the xz format and writes the compressed data to the output file.  The output file must end with the `.xz` extension.
func CompressFile(inpath string, outpath string, strategy ...CompressionStrategy) error {
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}
	return CompressFileWithProgress(inpath, outpath, nil, use_strategy)
}

// CompressFileWithProgress compresses a file using the xz format and writes the compressed data to the output file.  The output file must end with the `.xz` extension.
// Your progress callback function that you supply will be called with the number of bytes read and written to the output file.  This is useful for showing progress bars.
// The first 'uint64' is the number of bytes read from the input file, and the second 'uint64' is the number of bytes written to the output file.  From this you can calculate
// the percentage of the file that has been compressed, the estimated time remaining, etc.
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

	go func() {
		err := internal.CompressIn(readchan, writechan, int(use_strategy))
		if err != nil {
			fmt.Println("Error compressing data:", err)
		}
	}()
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

// CompressFileToMemory compresses a file using the xz format and returns the compressed data
// as a byte slice.  It can be handy for preparing uncompressed data for transmission over a network.
func CompressFileToMemory(path string, strategy ...CompressionStrategy) ([]byte, error) {

	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}

	f, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	go func() {
		err := internal.CompressIn(readchan, writechan, int(use_strategy))
		if err != nil {
			fmt.Println("Error compressing data:", err)
		}
	}()

	readfunc := func() {
		readbuf := make([]byte, internal.MAX_BUF_SIZE)

		for {
			bytes, err := f.Read(readbuf)
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

	membuffer := bytes.Buffer{}

	go readfunc()

	donewrite := make(chan bool, 1)
	go func() {
		for data := range writechan {
			membuffer.Write(data)
		}
		donewrite <- true
	}()
	<-donewrite

	return membuffer.Bytes(), nil
}

// CompressStream skips a call to io.Copy() by just compressing whatever stream you put in the
// input reader and writing it to the output writer.  If you hold the input stream open and keep writing to it,
// this call will block until you close the input stream.  This is useful for compressing data on the fly, such
// as the case with a logger stream that keeps the most recent events in RAM then shunts new entries off to
// a goroutine that's keeping a compressed version of it on disk.
// Note: Neiher CompressStream nor DecompressStream
// actually use XZReader or XZWriter.  They are just there for the sake of the ABI.
func CompressStream(input io.Reader, output io.Writer, strategy ...CompressionStrategy) error {
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}

	readchan := make(chan []byte, 1)
	writechan := make(chan []byte, 1)

	go func() {
		err := internal.CompressIn(readchan, writechan, int(use_strategy))
		if err != nil {
			fmt.Println("Error compressing data:", err)
		}
	}()

	readfunc := func() {
		readbuf := make([]byte, internal.MAX_BUF_SIZE)

		for {
			bytes, err := input.Read(readbuf)
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

	go readfunc()

	donewrite := make(chan bool, 1)
	go func() {
		for data := range writechan {
			output.Write(data)
		}
		donewrite <- true
	}()
	<-donewrite

	return nil
}
