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
// to pass 32/64-bit architecture to lzma since it must be set at compile time.
#cgo  386  amd64p32  arm  armbe  mips  mipsle  mips64p32  mips64p32le  ppc  riscv  s390  sparc CFLAGS: -DSIZEOF_SIZE_T=4
#cgo !386,!amd64p32,!arm,!armbe,!mips,!mipsle,!mips64p32,!mips64p32le,!ppc,!riscv,!s390,!sparc CFLAGS: -DSIZEOF_SIZE_T=8
// Tell C that we want the standard library
#cgo CFLAGS: -DHAVE_STDBOOL_H -DHAVE_STDINT_H -DHAVE_INTTYPES_H

#cgo LDFLAGS: -Linternal/lzma/src/liblzma -llzma

#include <stdlib.h>
#include <string.h>
#include "internal/lzma/src/liblzma/api/lzma.h"
#include "internal/common/sysdefs.h"

lzma_stream new_stream() {
	lzma_stream lz_stream = LZMA_STREAM_INIT;
	return lz_stream;
}
*/
import "C"
import (
	"unsafe"
)

const LZMA_CONCATENATED = 1

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

// Returns the number of bytes waiting to be fed into `liblzma`.  If this is non-zero,
// you can send more data to the stream (up to the number of bytes returned here).
func (s *lzmaStream) AvailableInputBytes() int {
	return int(s.cStream.avail_in)
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
	C.lzma_end(nil)
}

func Encoder(stream *lzmaStream, preset int) Return {
	return Return(C.lzma_easy_encoder(&stream.cStream, C.uint(preset), C.LZMA_CHECK_CRC64))
}

func Code(stream *lzmaStream, action Action) Return {
	return Return(C.lzma_code(&stream.cStream, C.lzma_action(action)))
}

const MAX_BUF_SIZE = 4096

func main() {

	str := "Hello, world!"
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	stream := createStream()
	defer stream.Close()
	Encoder(stream, 9)
	stream.SetInput([]byte("caltionyuotnwlgulwogh]<sientThis is a pretty big stringThis is a pretty big stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stri stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big stringThis is a pretty big string"))
	println("This is a pretty big string")
	Code(stream, Run)
	Code(stream, Finish)
	str = ""
	for stream.AvailableOutputBytes() > 0 {
		data := stream.Pop()
		if len(data) == 0 {
			break
		}
		str += string(data)
	}
	println(str)

}
