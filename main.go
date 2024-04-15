package main

/*
#cgo CFLAGS: -Iinternal/src/liblzma
#cgo CFLAGS: -Iinternal/src/liblzma/api
#cgo CFLAGS: -Iinternal/src/liblzma/common
#cgo CFLAGS: -Iinternal/src/liblzma/check
#cgo CFLAGS: -Iinternal/src/liblzma/delta
#cgo CFLAGS: -Iinternal/src/liblzma/lz
#cgo CFLAGS: -Iinternal/src/liblzma/lzma
#cgo CFLAGS: -Iinternal/src/liblzma/rangecoder
#cgo CFLAGS: -Iinternal/src/liblzma/simple
#cgo CFLAGS: -DHAVE_ENCODER_LZMA2 -DHAVE_DECODER_LZMA2
#cgo CFLAGS: -DHAVE_CHECK_CRC32 -DHAVE_CHECK_CRC64
#cgo CFLAGS: -DHAVE_MF_HC3 -DHAVE_MF_HC4 -DHAVE_MF_BT4

// This is taken from @jamespfennell's lzma-go project to get CGo a way
// to pass 32/64-bit architecture indicators to lzma since it must be set at compile time.
#cgo  386  amd64p32  arm  armbe  mips  mipsle  mips64p32  mips64p32le  ppc  riscv  s390  sparc CFLAGS: -DSIZEOF_SIZE_T=4
#cgo !386,!amd64p32,!arm,!armbe,!mips,!mipsle,!mips64p32,!mips64p32le,!ppc,!riscv,!s390,!sparc CFLAGS: -DSIZEOF_SIZE_T=8
// Tell C that we want the standard library
#cgo CFLAGS: -DHAVE_STDBOOL_H -DHAVE_STDINT_H -DHAVE_INTTYPES_H

// Switch on TUKLIB_OPTION_FAST_UNALIGNED_ACCESS to speed up the compression on x86 and x86_64 computers
#cgo 386 amd64 CFLAGS: -DTUKLIB_FAST_UNALIGNED_ACCESS

// Tell Cgo that we have liblzma source and turn on the C lzma macro
#cgo LDFLAGS: -Linternal/lzma/src/liblzma -llzma

#include <stdlib.h>
#include <string.h>
#include "internal/lzma/src/liblzma/api/lzma.h"
#include "internal/common/sysdefs.h"

// liblzma requires that the initialization of the stream be done with a C macro, which CGo cannot see.
// This function will not be called when this package init(), so it is safe to define it here.
lzma_mt multi_options = {
	.check = LZMA_CHECK_CRC64,
	.threads = 4,
};
lzma_mt get_multi_options() {
	return multi_options;
}
lzma_stream new_stream() {
	lzma_stream lz_stream = LZMA_STREAM_INIT;
	return lz_stream;
}
*/
import "C"
import (
	"os"
	"runtime"
	"unsafe"
)

// This is a type that is used to pass a buffer to C code. It is not safe to use from the top code.
type unsafeBuffer struct {
	start    *C.uint8_t
	length   C.size_t
	capacity C.size_t
}

// This function is used to grow the buffer. It is not safe to use from the top code.
func (b *unsafeBuffer) grow(size int) {
	if size <= int(b.capacity) {
		return
	}
	b.clear()
	b.start = (*C.uint8_t)(C.malloc(C.size_t(size)))
	b.length = 0
	b.capacity = C.size_t(size)
}

// This function is used to read from the buffer. It is not safe to use from the top code.
func (b *unsafeBuffer) read(length int) []byte {
	return C.GoBytes(unsafe.Pointer(b.start), C.int(length))
}

// This function is used to clear the buffer. It is not safe to use from the top code.
func (b *unsafeBuffer) clear() {
	if b.start != nil {
		C.free(unsafe.Pointer(b.start))
	}
	b.start = nil
	b.length = 0
	b.capacity = 0
}

func (b *unsafeBuffer) fill(data []byte) {
	if len(data) == 0 {
		b.length = 0
		return
	}
	b.grow(len(data))
	C.memcpy(unsafe.Pointer(b.start), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	b.length = C.size_t(len(data))
}

func (b *unsafeBuffer) toBytes(length int) []byte {
	if length == 0 { // If they don't want to read anything, send nothing.
		return []byte{}
	}
	return C.GoBytes(unsafe.Pointer(b.start), C.int(b.length))
}

// A wrapper around the internal C LZMA stream.
type lzmaStream struct {
	cStream C.lzma_stream
	input   unsafeBuffer
	output  unsafeBuffer
}

// Creates a new LZMA stream.
func createStream() *lzmaStream {
	stream := lzmaStream{
		cStream: C.new_stream(),
	}
	stream.output.grow(MAX_BUF_SIZE)
	stream.output.length = MAX_BUF_SIZE
	stream.cStream.next_out = stream.output.start
	stream.cStream.avail_out = stream.output.length
	return &stream
}

// Returns the number of bytes awaiting to be read out of the LZMA stream.  If this is non-zero,
// you can pop data from the stream (up to the number of bytes returned here).
func (s *lzmaStream) AvailableOutputBytes() int {
	return int(s.cStream.avail_out)
}

// Returns the number of bytes stacked on the input buffer (using SetInput) that
// are waiting to be processed by the LZMA stream.
func (s *lzmaStream) PendingInputBytes() int {
	return int(s.cStream.avail_in)
}

// Returns the total number of bytes that have been read from the input buffer.
func (s *lzmaStream) TotalInputBytes() int {
	return int(s.cStream.total_in)
}

// Returns the total number of bytes that have been written to the output buffer.
func (s *lzmaStream) TotalOutputBytes() int {
	return int(s.cStream.total_out)
}

// Pours data into the LZMA stream.  You have to make sure there is nothing in the buffer
// before calling or you this will overwrite data and cause corruption.
func (s *lzmaStream) SetInput(data []byte) {
	s.input.fill(data)
	s.cStream.next_in = s.input.start
	s.cStream.avail_in = s.input.length
}

// Pops out the data waiting to be read from the LZMA stream and clears the output buffer.
func (s *lzmaStream) Pop() []byte {
	buf := s.output.read(int(s.output.length - s.cStream.avail_out))
	s.cStream.next_out = s.output.start
	s.cStream.avail_out = s.output.length
	return buf
}

// Closes the internal LZMA stream and frees the memory.
func (s *lzmaStream) Close() {
	s.input.clear()
	s.output.clear()
	FreeLZMA(s)
}

// Frees the memory used by the LZMA stream.
func FreeLZMA(lzmaStream *lzmaStream) {
	C.lzma_end(&lzmaStream.cStream)
}

// The multi-threaded LZMA encoder.  Multi-threading doesn't do all that much for compression, but when
// you set compression to lower levels it can speed up the process.
func Encoder(stream *lzmaStream, preset int) Return {
	//Return(C.lzma_easy_encoder_mt(&stream.cStream, C.uint(preset), C.LZMA_CHECK_CRC64), C.uint(2))
	// Sets an LZMA stream up for an encoding job.
	options := C.get_multi_options()
	options.preset = C.uint(preset)
	options.threads = C.uint(runtime.NumCPU() / 2) // Use half the number of CPUs for encoding
	return Return(C.lzma_stream_encoder_mt(&stream.cStream, &options))
}

// Sets an LZMA stream up for a decoding job.
func Decoder(stream *lzmaStream, memlimit uint64) Return {
	return Return(C.lzma_stream_decoder(&stream.cStream, C.uint64_t(memlimit), C.LZMA_TELL_UNSUPPORTED_CHECK))
}

// Starts/Stops the LZMA stream encoding/decoding job.  This is a call-chain dependent function that
// requires Encoder or Decoder to be called first.
func EncodeDecodeJobAction(stream *lzmaStream, action Action) Return {
	return Return(C.lzma_code(&stream.cStream, C.lzma_action(action)))
}

const MAX_BUF_SIZE = 1328

func main() {

	str := "Hello, world!"
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	stream := createStream()
	defer stream.Close()
	encret := Encoder(stream, 6)
	println(encret.String())

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println("Current Directory:", currentDir)

	file, err := os.Open("test/canterbury-corpus/large/world192.txt")
	if err != nil {
		panic(err)
	}

	inbuffer := make([]byte, MAX_BUF_SIZE)
	outfile, _ := os.Create("world192.txt.xz")
	action := Run
	ret := Ok

	for {

		if action == Finish && ret == StreamEnd {
			break
		}

		bytesRead, err := file.Read(inbuffer)
		if err != nil {
			action = Finish
		}
		stream.SetInput(inbuffer[:bytesRead])

		ret = EncodeDecodeJobAction(stream, action)
		if ret != Ok && ret != StreamEnd {
			panic("Error in encoding/decoding job.")
		}

		print("\rBytes read: ", stream.TotalInputBytes(), " Bytes written: ", stream.TotalOutputBytes())
		outfile.Write(stream.Pop())

		if ret == StreamEnd {
			action = Finish
		}

	}
	println()
	outfile.Close()

}
