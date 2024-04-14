package main

/*
#cgo CFLAGS: -I./internal/src/liblzma/api
#include "lzma.h"
#include <stdlib.h>
#include <string.h>

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

func main() {

	str := "Hello, world!"
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

}
