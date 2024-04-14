# safexz
(in development)  A safe Go interface to liblzma, the "xz" compression library.

This is a Go package for compression in the xz / lzma format that provides a safer way to call `liblzma` for common use cases without the fear of type-safety issues and utilizes Go's goroutines to protect your project from [unforseen control hijacks](https://research.swtch.com/xz-timeline) from the 5.6.0 and 5.6.1 versions of `liblzma`

## compatibility with older platforms
`liblzma` is still one of the best compression algorithms for compacting data.  It has backportage to even OS/2 Warp from the mid-1990s and there is an implementation for nearly every 32-bit-or-better processor.

As Go itself is much newer, I am to maintain compatibility `$TODO: backportage tests :-)`.

## limited but easy interface
Rather than expose the call-chain dependencies of `lzma` directly to you, a simpler Go interface is provided for your integration projects.  The interface breaks down your calls into simpler `chan` transfers to `lzma`.

## safexz wraps the lzma library in isolated go routines
This project's goal is to abide by a chaos/containment theory which believes _unless you wrote the C code, tread carefully_.

To increase the difficulty of a user with a suspicious C library gaining control over your own go program that uses a C library, `safexz`:

- Hides C call access and data structures as internal packages
- No call is made to a C library without boxing it inside 2 layers of goroutines
- All communication of data to and from the C library's functions must take place from a separate goroutine and data may only pass via channels using packed types (fields without pointers) wherever possible.
- Go pointers into `lzma` are not possible
- No one using this library would ever have to manage `liblzma`'s internal state, nor will `liblzma` be able to "see" your data types of your own code except for the []byte stream of data to be compressed/decompressed.

More detail about this technique is [over here](https://gist.github.com/christoofar/880b4bcf3018f4681bb71bfdf1c16a6a).