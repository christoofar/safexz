# Development Log

## Apr 19 2024

So, I finally figured out why I was getting such non-deterministic results
![image](https://github.com/christoofar/safexz/assets/5059144/f4136c9c-3742-4262-9a26-03eacd338ac0)

`readbuf` is dirty after the read.  I thought that this was a simple reusable type but it's got some distinct behavior when used with `file.Read()` because of the syscall that occurs.   To fix this, I pull out what was read into a clean byte slice and send that into the `chan` for processing.

That results in a clean byte-for-byte accounted-for `diff`:

![image](https://github.com/christoofar/safexz/assets/5059144/98292003-274d-4375-bb88-6e413aa6726a)

Now I can move on to the multi-threaded decompression.

## Apr 18 2024

So, all options using Simple (single-threading) produce a good result.  The multi-threading ones do not.  Probably another signal that I need to pick up from `liblzma.so` to know that all the threads underneath in the innermost goroutine have completed.  I'll be hunting around for some multithread examples in C to see if the calling pattern is different.

*Eureka!* I found what was wrong with the multithreaded compress options.  Turns out that I goofed and did not check `.avail_in` on the stream before pushing data.  Apparently this issue doesn't turn up frequently enough for me to see it in single-threaded mode but it will come up in multi-threaded.   `.avail_in` tells you that there are bytes waiting to be drained into the memory area where lzma is working.   You can try to fill up to `MAX_BUFFER_SIZE` but it's easier to wait for it to clear to zero in a cycle and on the next cycle it's likely for the drain to occur.  If you just set all the bytes for the cycle then the bytes waiting to be drained will be destroyed.

Now, testing again with compressing the Debian12 DVD I still get two different compression result sizes.
![image](https://github.com/christoofar/safexz/assets/5059144/ef8fdbec-957e-4d2c-8023-c338f656b996)

But this time the byte counts matches the origin.
![image](https://github.com/christoofar/safexz/assets/5059144/4120ddad-c460-4481-8841-a3909c0d1b82)

Is LZMA still not deterministic?  That's wild.  Can't be right.  Let's check that all the bytes exactly match.
![image](https://github.com/christoofar/safexz/assets/5059144/cc43aa82-99b2-47d5-ab38-31d8063906da)

You know... I think maybe `MAX_BUFFER_SIZE` really should be down in the low count, like `1024`.  Anything higher might break some limits or plausables. 

## Apr 17 2024

- Completed encoding pathway using `*<-chan` and `*chan<-` streaming paths to `liblzma.so`
- Set up Go interface stubs.  First one to be wired up is `CompressFileWithProgress` in `compress.go`
- Set up a pre-defined matrix of compression strategies to simplify the options in `liblzmna.so`, these are:

|Strategy  | Threads  |  `liblzma.so` level |
|----------|----------|---------------------|
|CompressionSimple      | 1 |  4  |
|CompressionSimpleFast  | 1 |  2  |
|CompressionSimpleBetter| 1 |  7  |
|CompressionSimpleMax   | 1 |  9  |
|CompressionMulti*       | ½ vCPUs | 4 |
|CompressionMultiFast   | ½ vCPUs | 2 |
|CompressionMultiBetter | ½ vCPUs | 7 |
|CompressionMultiMax    | ½ vCPUs | 9 |
|CompressionFullPower         | All vCPUs | 4 |
|CompressionFullPowerFast     | All vCPUs | 2 |
|CompressionFullPowerBetter   | All vCPUs | 7 |
|CompressionFullPowerMax      | All vCPUs | 9 |

`* If only 2 cores are available, the Multi option will use both cores just the same as asking for the FullPower strategy.   On a single-core machine, all the options will default to single thread.`

I chose `CompressionMulti` as the default option for unspecified, as it balances for a multi-core environment which most people are running in containers but stays conservative on the memory so only 300MB max is expected to go into reservation when passing in a big file.   If you own a beast of a ThreadRipper you obviously are going to reach for `CompressionFullPowerMax`, while the default setting will still work fine on older Raspberry PI boards.

- To check for memory leaks I wrote a simple loop to bang away at file compression using my copy of the Debian12 installation DVD.  I was surprised that `xz` is a bit non-deterministic in compression results, which is not an expected outcome.   Pass 3 compressed to 2,640,800,208 bytes while Pass 10 compressed to 2,600,824,772.   I checked the progress report calls to make sure I wasn't missing a call to report progress at the tail end of compression.
![image](https://github.com/christoofar/safexz/assets/5059144/0ff252d7-41b6-4d9f-8afc-d781d095d6d4)

Let's decompress the result of Pass 3...
![image](https://github.com/christoofar/safexz/assets/5059144/334658d3-1134-4f14-999b-c68d32c0d598)

Looks like I have a bug somewhere in my stream mechanism.   I'll need to make sure the input bytes are all getting into the encoder and that I am not prematurely closing out of the encoding cycle.

- At least the memory leak checks are passing.
