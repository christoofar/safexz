// Exposes the internal functions of lzma as a construction of calls requiring streaming
// input and output channels.
package internal

func CompressIn(in *chan []byte, out *chan []byte) {

	// The reason for the nested go routines is to isolate the unsafeBuffer
	go func(input *chan []byte, output *chan []byte) {
		go func(receive <-chan []byte, sender chan<- []byte) {
			compressChanStream(&receive, &sender)
		}(*input, *output)
	}(in, out)

}